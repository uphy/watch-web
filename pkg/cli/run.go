package cli

import (
	"fmt"

	"github.com/urfave/cli"
)

func (c *CLI) run() cli.Command {
	return cli.Command{
		Name: "run",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name: "a,all",
			},
			cli.BoolFlag{
				Name: "t,test-action",
			},
		},
		Action: func(ctx *cli.Context) error {
			all := ctx.Bool("all")
			testAction := ctx.Bool("test-action")
			exe := c.executor
			if all {
				exe.CheckAll()
			} else {
				for _, id := range ctx.Args() {
					job := exe.Job(id)
					result, err := exe.Check(job)
					if err != nil {
						continue
					}
					fmt.Println("[Result]")
					fmt.Println(result.Current.YAML())

					fmt.Println("[Diff]")
					diff := result.Diff()
					fmt.Println(diff.YAML())

					if testAction {
						if err := exe.DoActions(job, result); err != nil {
							fmt.Println("failed on action: ", err)
						}
					}
				}
			}
			return nil
		},
	}
}
