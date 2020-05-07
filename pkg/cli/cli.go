package cli

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/uphy/watch-web/pkg/config"
	"github.com/uphy/watch-web/pkg/watch"
	"github.com/urfave/cli"
)

type CLI struct {
	app      *cli.App
	executor *watch.Executor
	log      *logrus.Logger
}

func New(log *logrus.Logger) *CLI {
	app := cli.NewApp()
	app.Name = "watch-web"
	app.Usage = "Watch web updated"
	c := &CLI{app, nil, log}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "c,config",
			Value: "config.yml",
		},
		cli.BoolFlag{
			Name: "debug",
		},
		cli.BoolFlag{
			Name: "info",
		},
	}
	app.Before = func(ctx *cli.Context) error {
		// set log level based on flags
		if ctx.Bool("debug") {
			c.log.SetLevel(logrus.DebugLevel)
		} else if ctx.Bool("info") {
			c.log.SetLevel(logrus.InfoLevel)
		} else {
			c.log.SetLevel(logrus.WarnLevel)
		}
		// load config file
		configFile := ctx.String("config")
		e, err := config.LoadAndCreate(c.log, configFile)
		if err != nil {
			return fmt.Errorf("failed to load config file: %w", err)
		}
		c.executor = e
		log.WithFields(logrus.Fields{
			"file": configFile,
		}).Debug("Config file loaded.")
		return nil
	}
	app.Commands = []cli.Command{
		c.start(),
		c.run(),
		c.list(),
	}
	return c
}

func (c *CLI) Run(args []string) error {
	return c.app.Run(args)
}
