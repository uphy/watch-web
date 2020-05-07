package config

import (
	"github.com/uphy/watch-web/pkg/domain"
)

type (
	ActionConfig struct {
		Slack   *SlackActionConfig   `json:"slack,omitempty"`
		Console *ConsoleActionConfig `json:"console"`
	}
	SlackActionConfig struct {
		URL   domain.TemplateString `json:"url"`
		Debug bool                  `json:"debug"`
	}
	ConsoleActionConfig struct {
	}
)
