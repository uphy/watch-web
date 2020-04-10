package check

import (
	"os/exec"
)

type (
	CommandSource struct {
		Command   string `json:"command"`
		LabelText string `json:"label"`
	}
)

func NewCommandSource(command, label string) *CommandSource {
	return &CommandSource{
		Command:   command,
		LabelText: label,
	}
}

func (c *CommandSource) Label() string {
	return c.LabelText
}

func (c *CommandSource) Fetch() (string, error) {
	cmd := exec.Command("sh", "-c", c.Command)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(b), nil
}
