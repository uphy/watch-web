package config

import (
	"github.com/uphy/watch-web/pkg/domain/template"
)

type (
	ActionConfig struct {
		SlackWebhook *SlackWebhookActionConfig `json:"slack_webhook,omitempty"`
		Console      *ConsoleActionConfig      `json:"console"`
	}
	SlackWebhookActionConfig struct {
		URL   template.TemplateString `json:"url"`
		Debug bool                    `json:"debug"`
	}
	ConsoleActionConfig struct {
	}
)
