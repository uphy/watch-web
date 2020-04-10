package check

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "PATH=") {
			dir, err := os.Getwd()
			if err != nil {
				return "", err
			}
			paths := []string{filepath.Join(dir, "scripts")}
			cmd.Env = append(cmd.Env, "PATH="+strings.Join(paths, ":")+":"+env[5:])
		} else {
			cmd.Env = append(cmd.Env, env)
		}
	}
	//cmd.Env = []string{"PATH=/Users/ishikura/dev/github.com/uphy/watch-web/server/scripts:" + os.Getenv("PATH")}
	// fmt.Println(cmd.Env)
	// cmd.Env = append([]string{"PATH=" + path}, os.Environ()...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(b), nil
}
