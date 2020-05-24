package config

import (
	"github.com/uphy/watch-web/pkg/domain/template"
)

const (
	SlackActionThreadPerJob  = "job"
	SlackActionThreadPerItem = "item"
)

type (
	ActionConfig struct {
		SlackWebhook *SlackWebhookActionConfig `json:"slack_webhook,omitempty"`
		SlackBot     *SlackBotActionConfig     `json:"slack,omitempty"`
		Console      *ConsoleActionConfig      `json:"console"`
	}
	SlackWebhookActionConfig struct {
		URL   template.TemplateString `json:"url"`
		Debug bool                    `json:"debug"`
	}
	SlackBotActionConfig struct {
		Token     template.TemplateString  `json:"token"`
		Channel   template.TemplateString  `json:"channel"`
		ThreadPer *template.TemplateString `json:"thread_per"`
	}
	ConsoleActionConfig struct {
	}
)
