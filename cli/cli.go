package cli

import (
	"fmt"

	"github.com/uphy/watch-web/config"
	"github.com/urfave/cli"
)

type CLI struct {
	app    *cli.App
	config *config.Config
}

func New() *CLI {
	app := cli.NewApp()
	app.Name = "watch-web"
	app.Usage = "Watch web updated"
	c := &CLI{app, nil}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "c,config",
			Value: "config.yml",
		},
	}
	app.Before = func(ctx *cli.Context) error {
		configFile := ctx.String("config")
		conf, err := config.LoadConfigFile(configFile)
		if err != nil {
			return fmt.Errorf("failed to load config file: %w", err)
		}
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
