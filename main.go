package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/uphy/watch-web/check"
	"github.com/uphy/watch-web/config"
	"github.com/uphy/watch-web/resources"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	conf, err := config.LoadConfigFile("./config.yml")
	if err != nil {
		return err
	}
	exe, err := conf.NewExecutor()
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
			if !ptn.Match([]byte(j.Name)) {
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
	e.GET("/*", echo.WrapHandler(http.FileServer(resources.HttpStatic)))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return e.Start(fmt.Sprintf(":%s", port))
}
