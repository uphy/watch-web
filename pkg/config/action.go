package config

import (
	"errors"
	"time"

	"github.com/uphy/watch-web/pkg/domain"
	"github.com/uphy/watch-web/pkg/domain/template"
	"github.com/uphy/watch-web/pkg/watch/store"
)

const (
	SlackActionThreadPerJob  = "job"
	SlackActionThreadPerItem = "item"
)

type (
	ActionConfig struct {
		SlackWebhook *SlackWebhookActionConfig `json:"slack_webhook,omitempty"`
		SlackBot     *SlackBotActionConfig     `json:"slack_bot,omitempty"`
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
		Debug     bool                     `json:"debug"`
	}
	ConsoleActionConfig struct {
	}

	slackBotActionRepository struct {
		store     domain.Store
		threadPer string
	}
)

func (s *SlackBotActionConfig) newRepository(ctx *template.TemplateContext, store domain.Store) (*slackBotActionRepository, error) {
	var threadPer string
	if s.ThreadPer != nil {
		var err error
		threadPer, err = s.ThreadPer.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		if threadPer != SlackActionThreadPerItem && threadPer != SlackActionThreadPerJob {
			return nil, errors.New("unsupported thread_per param: " + threadPer)
		}
	} else {
		threadPer = SlackActionThreadPerJob
	}
	return &slackBotActionRepository{store, threadPer}, nil
}

func (s *slackBotActionRepository) PutSlackThreadTS(jobID string, itemID string, threadTS string) error {
	return s.store.SetTemp(s.key(jobID, itemID), threadTS, time.Hour*24*30)
}

func (s *slackBotActionRepository) GetSlackThreadTS(jobID, itemID string) (*string, error) {
	v, err := s.store.Get(s.key(jobID, itemID))
	if err == store.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (s *slackBotActionRepository) key(jobID string, itemID string) string {
	if s.threadPer == SlackActionThreadPerJob {
		return "slackts:" + jobID
	} else {
		return "slackts:" + jobID + ":" + itemID
	}
}
