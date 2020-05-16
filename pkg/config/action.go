package config

import (
	"github.com/uphy/watch-web/pkg/domain/template"
)

type (
	ActionConfig struct {
		Slack   *SlackActionConfig   `json:"slack,omitempty"`
		Console *ConsoleActionConfig `json:"console"`
	}
	SlackActionConfig struct {
		URL   template.TemplateString `json:"url"`
		Debug bool                    `json:"debug"`
	}
	ConsoleActionConfig struct {
	}
)
