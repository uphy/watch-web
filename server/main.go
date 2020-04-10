package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/uphy/watch-web/server/check"
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
		for _, j := range exe.Jobs {
			j.Check()
		}
		return ctx.JSON(200, exe.Jobs)
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
		return ctx.NoContent(200)
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
	e.Static("/", "./static")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return e.Start(fmt.Sprintf(":%s", port))
}
