package check

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func (c *ShellSource) Fetch() (string, error) {
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
	b, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(b), nil
}
