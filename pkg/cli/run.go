package cli

import (
	"fmt"

	"github.com/uphy/watch-web/pkg/value"
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
			exe, err := c.newExecutor()
			if err != nil {
				return fmt.Errorf("failed to create executor: %w", err)
			}
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
					fmt.Println(result.Current)

					fmt.Println("[Diff]")
					if result.Previous == "" {
						var prev string
						switch result.ValueType {
						case value.ValueTypeString:
							prev = ""
						case value.ValueTypeJSONArray:
							prev = "[]"
						case value.ValueTypeJSONObject:
							prev = "{}"
						}
						result.Previous = prev
					}
					diff, err := result.Diff()
					if err != nil {
						fmt.Println("failed on diff: ", err)
					} else {
						fmt.Println(diff)
					}

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
