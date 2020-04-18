package check

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/uphy/watch-web/pkg/value"
)

type (
	ShellSource struct {
		Command   string
		ValueType value.ValueType
	}
)

func NewShellSource(command string, valueType value.ValueType) *ShellSource {
	return &ShellSource{
		Command:   command,
		ValueType: valueType,
	}
}

func (c *ShellSource) Fetch(ctx *JobContext) (value.Value, error) {
	ctx.Log.WithField("command", c.Command).Debug("Run shell command.")
	cmd := exec.Command("sh", "-c", c.Command)
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "PATH=") {
			dir, err := os.Getwd()
			if err != nil {
				return nil, err
			}
			paths := []string{filepath.Join(dir, "scripts")}
			cmd.Env = append(cmd.Env, "PATH="+strings.Join(paths, ":")+":"+env[5:])
		} else {
			cmd.Env = append(cmd.Env, env)
		}
	}
	b, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return value.ConvertAs(string(b), c.ValueType)
}

func (c *ShellSource) String() string {
	return fmt.Sprintf("Shell[command=%s]", c.Command)
}
