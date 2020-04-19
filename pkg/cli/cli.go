package cli

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/uphy/watch-web/pkg/watch"
	"github.com/uphy/watch-web/pkg/config"
	"github.com/urfave/cli"
)

type CLI struct {
	app    *cli.App
	config *config.Config
	log    *logrus.Logger
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
		conf, err := config.LoadConfigFile(configFile)
		if err != nil {
			return fmt.Errorf("failed to load config file: %w", err)
		}
		log.WithFields(logrus.Fields{
			"file": configFile,
		}).Debug("Config file loaded.")
		c.config = conf
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

func (c *CLI) newExecutor() (*watch.Executor, error) {
	return c.config.NewExecutor(c.log)
}
