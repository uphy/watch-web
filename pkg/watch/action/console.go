package action

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
	updates := res.Diff()

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("%s (%s)\n", res.Label, res.JobID)
	fmt.Println("--------------------------------------------------------------------------------")
	if updates.Changes() {
		fmt.Println("Changes:")
		fmt.Println(updates)
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
