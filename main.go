package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/uphy/watch-web/check"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	config, err := check.LoadConfigFile("./config.yml")
	if err != nil {
		return err
	}
	exe, err := config.NewExecutor()
	if err != nil {
		return err
	}
	exe.Start()

	e := echo.New()
	e.GET("/api/jobs", func(ctx echo.Context) error {
		jobs := []string{}
		for _, j := range exe.Jobs {
			jobs = append(jobs, j.Name)
		}
		return ctx.JSON(200, jobs)
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
		job.Check()
		return ctx.Redirect(301, fmt.Sprintf("/api/jobs/%s", name))
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
	return e.Start(":8080")
}
