package config

import (
	"errors"
	"fmt"

	"github.com/uphy/watch-web/pkg/check"
	"github.com/uphy/watch-web/pkg/config/template"
)

type (
	ActionConfig struct {
		Slack      *SlackActionConfig      `json:"slack,omitempty"`
		LINENotify *LINENotifyActionConfig `json:"line_notify,omitempty"`
		Console    *ConsoleActionConfig    `json:"console"`
	}
	SlackActionConfig struct {
		URL   template.TemplateString `json:"url"`
		Debug bool                    `json:"debug"`
	}
	LINENotifyActionConfig struct {
		AccessToken template.TemplateString `json:"access_token"`
	}
	ConsoleActionConfig struct {
	}
)

func (a *ActionConfig) Action(ctx *template.TemplateContext) (check.Action, error) {
	if a.Slack != nil {
		return a.Slack.Action(ctx)
	}
	if a.LINENotify != nil {
		return a.LINENotify.Action(ctx)
	}
	fmt.Println(a.Console)
	if a.Console != nil {
		return check.NewConsoleAction(), nil
	}
	return nil, errors.New("no action defined")
}

func (s *SlackActionConfig) Action(ctx *template.TemplateContext) (check.Action, error) {
	url, err := s.URL.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return check.NewSlackAction(url, s.Debug), nil
}

func (s *LINENotifyActionConfig) Action(ctx *template.TemplateContext) (check.Action, error) {
	token, err := s.AccessToken.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return check.NewLINENotifyAction(token), nil
}
