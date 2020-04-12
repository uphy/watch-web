package cli

import (
	"github.com/urfave/cli"
)

func (c *CLI) run() cli.Command {
	return cli.Command{
		Name: "run",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name: "a,all",
			},
		},
		Action: func(ctx *cli.Context) error {
			all := ctx.Bool("all")
			exe, err := c.config.NewExecutor()
			if err != nil {
				return err
			}
			if all {
				for _, job := range exe.Jobs {
					job.Check()
				}
			} else {
				for _, id := range ctx.Args() {
					exe.Job(id).Check()
				}
			}
			return nil
		},
	}
}
