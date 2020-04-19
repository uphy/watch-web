package actions

import (
	"fmt"

	"github.com/uphy/watch-web/pkg/domain"
)

type (
	ConsoleAction struct {
	}
)

func NewConsoleAction() *ConsoleAction {
	return &ConsoleAction{}
}

func (s *ConsoleAction) Run(ctx *domain.JobContext, res *domain.Result) error {
	changes, err := res.Diff()
	if err != nil {
		return err
	}

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("%s (%s)\n", res.Label, res.JobID)
	fmt.Println("--------------------------------------------------------------------------------")
	if changes.Changed() {
		fmt.Println("Changes:")
		fmt.Println(changes.String())
		fmt.Println()
		fmt.Println("Previous:")
		fmt.Println(res.Previous)
		fmt.Println("Current:")
		fmt.Println(res.Current)
	} else {
		fmt.Println("Not Changed")
		fmt.Println("Current:")
		fmt.Println(res.Current)
	}
	return nil
}