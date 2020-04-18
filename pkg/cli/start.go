package cli

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/uphy/watch-web/pkg/check"
	"github.com/uphy/watch-web/pkg/resources"
	"github.com/urfave/cli"
)

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
				jobs := []check.Job{}
				for _, j := range exe.Jobs {
					job := *j
					job.Previous = nil
					jobs = append(jobs, job)
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
					if !ptn.Match([]byte(j.ID)) {
						continue
					}
					j.Check()
					checked = append(checked, j)
				}
				return ctx.JSON(200, checked)
			})
			e.GET("/api/jobs/:name", func(ctx echo.Context) error {
				name := ctx.Param("name")
				job := exe.Job(name)
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
				result, err := job.Check()
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
				if err := job.TestActions(); err != nil {
					return echo.NewHTTPError(500, err)
				}
				return ctx.NoContent(200)
			})
			e.GET("/*", echo.WrapHandler(http.FileServer(resources.HttpStatic)))
			return e.Start(fmt.Sprintf(":%d", port))
		},
	}
}
