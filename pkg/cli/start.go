package cli

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/uphy/watch-web/pkg/check"
	"github.com/uphy/watch-web/pkg/resources"
	"github.com/urfave/cli"
)

type (
	JobDTO struct {
		ID     string       `json:"id"`
		Link   string       `json:"link,omitempty"`
		Label  string       `json:"label,omitempty"`
		Status check.Status `json:"status"`
		Error  *string      `json:"error,omitempty"`
		Last   *time.Time   `json:"last,omitempty"`
		Count  int          `json:"count"`
	}
	JobDetailDTO struct {
		*JobDTO
		Previous *string `json:"previous,omitempty"`
	}
)

func getJob(exe *check.Executor, jobID string) (*JobDTO, error) {
	job := exe.Jobs[jobID]
	status, err := exe.GetJobStatus(jobID)
	if err != nil {
		return nil, err
	}
	if status == nil {
		status = &check.JobStatus{
			Status: check.StatusOK,
		}
	}
	return &JobDTO{
		ID:     job.ID(),
		Link:   job.Info.Link,
		Label:  job.Info.Label,
		Status: status.Status,
		Error:  status.Error,
		Last:   status.Last,
		Count:  status.Count,
	}, nil
}

func getJobs(exe *check.Executor) ([]JobDTO, error) {
	jobs := []JobDTO{}
	for jobID := range exe.Jobs {
		job, err := getJob(exe, jobID)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, *job)
	}
	return jobs, nil
}

func getJobDetail(exe *check.Executor, jobID string) (*JobDetailDTO, error) {
	job, err := getJob(exe, jobID)
	if err != nil {
		return nil, err
	}
	value, err := exe.GetJobValue(jobID)
	if err != nil {
		return nil, err
	}
	return &JobDetailDTO{job, value}, nil
}

func (c *CLI) start() cli.Command {
	return cli.Command{
		Name: "start",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name: "a,api",
			},
			cli.BoolFlag{
				Name: "ns,no-schedule",
			},
			cli.IntFlag{
				Name:   "p,port",
				Value:  8080,
				EnvVar: "PORT",
			},
		},
		Action: func(ctx *cli.Context) error {
			enableAPI := ctx.Bool("api")
			enableSchedule := !ctx.Bool("no-schedule")
			exe, err := c.newExecutor()
			if err != nil {
				return err
			}

			if !enableAPI {
				if enableSchedule {
					exe.Run()
				}
				return nil
			}
			if enableSchedule {
				go exe.Run()
			}

			port := ctx.Int("port")
			e := echo.New()
			e.Pre(middleware.RemoveTrailingSlash())
			e.Use(middleware.Logger())
			e.GET("/api/jobs", func(ctx echo.Context) error {
				jobs, err := getJobs(exe)
				if err != nil {
					return err
				}
				return ctx.JSON(200, jobs)
			})
			e.POST("/api/jobs/check-all", func(ctx echo.Context) error {
				v, err := ctx.FormParams()
				if err != nil {
					return err
				}
				p := v.Get("pattern")
				ptn, err := regexp.Compile(p)
				if err != nil {
					return err
				}
				checked := make([]*check.Job, 0)
				for _, j := range exe.Jobs {
					if !ptn.Match([]byte(j.ID())) {
						continue
					}
					exe.Check(j)
					checked = append(checked, j)
				}
				return ctx.JSON(200, checked)
			})
			e.GET("/api/jobs/:jobID", func(ctx echo.Context) error {
				jobID := ctx.Param("jobID")
				job, err := getJobDetail(exe, jobID)
				if err != nil {
					return err
				}
				if job == nil {
					return echo.NewHTTPError(404, "specified job is not exist")
				}
				return ctx.JSON(200, job)
			})
			e.POST("/api/jobs/:name/check", func(ctx echo.Context) error {
				name := ctx.Param("name")
				job := exe.Job(name)
				if job == nil {
					return echo.NewHTTPError(404, "specified job is not exist")
				}
				result, err := exe.Check(job)
				if err != nil {
					return echo.NewHTTPError(500, "failed to check: "+err.Error())
				}
				return ctx.JSON(200, result)
			})
			e.POST("/api/jobs/:name/test-actions", func(ctx echo.Context) error {
				name := ctx.Param("name")
				job := exe.Job(name)
				if job == nil {
					return echo.NewHTTPError(404, "specified job is not exist")
				}
				if err := exe.TestActions(job); err != nil {
					return echo.NewHTTPError(500, err)
				}
				return ctx.NoContent(200)
			})
			e.GET("/*", echo.WrapHandler(http.FileServer(resources.HttpStatic)))
			return e.Start(fmt.Sprintf(":%d", port))
		},
	}
}
