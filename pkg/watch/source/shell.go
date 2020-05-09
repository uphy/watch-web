package source

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/uphy/watch-web/pkg/domain"
)

type (
	ShellSource struct {
		Command string
	}
)

func NewShellSource(command string) *ShellSource {
	return &ShellSource{
		Command: command,
	}
}

func (c *ShellSource) Fetch(ctx *domain.JobContext) (domain.Value, error) {
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
		return nil, fmt.Errorf("failed to execute shell command: command=%s, err=%v", c.Command, err)
	}
	return domain.NewStringValue(string(b)), nil
}

func (c *ShellSource) String() string {
	return fmt.Sprintf("Shell[command=%s]", c.Command)
}
