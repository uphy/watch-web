package config

import (
	"errors"
	"fmt"

	"github.com/uphy/watch-web/pkg/domain"
	actions "github.com/uphy/watch-web/pkg/watch/action"
)

type (
	ActionConfig struct {
		Slack      *SlackActionConfig      `json:"slack,omitempty"`
		LINENotify *LINENotifyActionConfig `json:"line_notify,omitempty"`
		Console    *ConsoleActionConfig    `json:"console"`
	}
	SlackActionConfig struct {
		URL   domain.TemplateString `json:"url"`
		Debug bool                  `json:"debug"`
	}
	LINENotifyActionConfig struct {
		AccessToken domain.TemplateString `json:"access_token"`
	}
	ConsoleActionConfig struct {
	}
)

func (a *ActionConfig) Action(ctx *domain.TemplateContext) (domain.Action, error) {
	if a.Slack != nil {
		return a.Slack.Action(ctx)
	}
	if a.LINENotify != nil {
		return a.LINENotify.Action(ctx)
	}
	fmt.Println(a.Console)
	if a.Console != nil {
		return actions.NewConsoleAction(), nil
	}
	return nil, errors.New("no action defined")
}

func (s *SlackActionConfig) Action(ctx *domain.TemplateContext) (domain.Action, error) {
	url, err := s.URL.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return actions.NewSlackAction(url, s.Debug), nil
}

func (s *LINENotifyActionConfig) Action(ctx *domain.TemplateContext) (domain.Action, error) {
	token, err := s.AccessToken.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return actions.NewLINENotifyAction(token), nil
}
