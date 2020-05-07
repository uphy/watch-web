package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/urfave/cli"
)

func (c *CLI) list() cli.Command {
	return cli.Command{
		Name: "list",
		Action: func(ctx *cli.Context) error {
			exe := c.executor
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
			for _, job := range exe.Jobs {
				fmt.Fprintf(w, "%s\t%s\t%s\n", job.ID(), job.Info.Label, job.Info.Link)
			}
			w.Flush()
			return nil
		},
	}
}
