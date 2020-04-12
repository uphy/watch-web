package config

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/go-redis/redis/v7"
	"github.com/uphy/watch-web/check"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
)

func (c *Config) NewExecutor() (*check.Executor, error) {
	ctx := NewRootContext()
	store, err := newStore(ctx, c.Store)
	if err != nil {
		return nil, err
	}
	e := check.NewExecutor(store)
	if c.InitialRun != nil {
		e.InitialRun = *c.InitialRun
	}
	for _, jobConfig := range c.Jobs {
		if len(jobConfig.WithItems) > 0 {
			for itemIndex, item := range jobConfig.WithItems {
				evaluatedItem, err := evaluateItemAsTemplate(ctx, item)
				if err != nil {
					return nil, err
				}
				ctx.PushScope()
				ctx.Set("itemIndex", itemIndex)
				ctx.Set("item", evaluatedItem)
				addJob(ctx, e, &jobConfig)
				ctx.PopScope()
			}
		} else {
			addJob(ctx, e, &jobConfig)
		}
	}
	return e, nil
}

func evaluateItemAsTemplate(ctx *TemplateContext, v interface{}) (interface{}, error) {
	m, ok := v.(map[string]interface{})
	if ok {
		evaluated := make(map[string]interface{})
		for key, value := range m {
			ekey, err := TemplateString(key).Evaluate(ctx)
			if err != nil {
				return nil, err
			}
			evalue, err := evaluateItemAsTemplate(ctx, value)
			if err != nil {
				return nil, err
			}
			evaluated[ekey] = evalue
		}
		return evaluated, nil
	}
	e, err := TemplateString(fmt.Sprint(v)).Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func addJob(ctx *TemplateContext, e *check.Executor, jobConfig *JobConfig) error {
	source, err := jobConfig.Source.Source(ctx)
	if err != nil {
		return err
	}
	actions := []check.Action{}
	for _, actionConfig := range jobConfig.Actions {
		action, err := actionConfig.Action(ctx)
		if err != nil {
			return err
		}
		actions = append(actions, action)
	}
	id, err := jobConfig.ID.Evaluate(ctx)
	if err != nil {
		return err
	}
	schedule, err := jobConfig.Schedule.Evaluate(ctx)
	if err != nil {
		return err
	}
	label, err := jobConfig.Label.Evaluate(ctx)
	if err != nil {
		return err
	}
	link, err := jobConfig.Link.Evaluate(ctx)
	if err != nil {
		return err
	}
	if err := e.AddJob(id, schedule, label, link, source, actions...); err != nil {
		return err
	}
	return nil
}

func (s *SourceConfig) Source(ctx *TemplateContext) (check.Source, error) {
	if s.DOM != nil {
		return s.DOM.Source(ctx)
	}
	if s.Shell != nil {
		return s.Shell.Source(ctx)
	}
	return nil, errors.New("no source defined")
}

func (d *DOMSourceConfig) Source(ctx *TemplateContext) (check.Source, error) {
	var encoding encoding.Encoding
	if d.Encoding != nil {
		enc, err := d.Encoding.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		switch enc {
		case "Shift_JIS", "sjis":
			encoding = japanese.ShiftJIS
		default:
			return nil, errors.New("unsupported encoding: " + enc)
		}
	}
	u, err := d.URL.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	s, err := d.Selector.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	source := check.NewDOMSource(u, s, encoding)
	return source, nil
}

func (d *ShellSourceConfig) Source(ctx *TemplateContext) (check.Source, error) {
	command, err := d.Command.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return check.NewShellSource(command), nil
}
func (a *ActionConfig) Action(ctx *TemplateContext) (check.Action, error) {
	if a.Slack != nil {
		return a.Slack.Action(ctx)
	}
	if a.LINENotify != nil {
		return a.LINENotify.Action(ctx)
	}
	return nil, errors.New("no action defined")
}

func (s *SlackActionConfig) Action(ctx *TemplateContext) (check.Action, error) {
	url, err := s.URL.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return check.NewSlackAction(url), nil
}

func (s *LINENotifyActionConfig) Action(ctx *TemplateContext) (check.Action, error) {
	token, err := s.AccessToken.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return check.NewLINENotifyAction(token), nil
}

func newStore(ctx *TemplateContext, config *StoreConfig) (check.Store, error) {
	if config != nil && config.Redis != nil {
		password := ""
		addr := ""
		if config.Redis.Address != nil {
			if config.Redis.Password != nil {
				p, err := config.Redis.Password.Evaluate(ctx)
				if err != nil {
					return nil, err
				}
				password = p
			}
			a, err := config.Redis.Address.Evaluate(ctx)
			if err != nil {
				return nil, err
			}
			addr = a
		} else if config.Redis.RedisToGo != nil {
			r, err := config.Redis.RedisToGo.Evaluate(ctx)
			if err != nil {
				return nil, err
			}
			addr, password, err = parseRedisToGoURL(r)
			if err != nil {
				log.Println(err)
				return nil, err
			}
		}
		if addr != "" {
			client := redis.NewClient(&redis.Options{
				Addr:     addr,
				Password: password,
			})
			return check.NewRedisStore(client), nil
		}
	}
	return &check.NullStore{}, nil
}

func parseRedisToGoURL(redisToGo string) (addr string, password string, err error) {
	var redisInfo *url.URL
	redisInfo, err = url.Parse(redisToGo)
	if err != nil {
		return
	}

	addr = redisInfo.Host
	if redisInfo.User != nil {
		password, _ = redisInfo.User.Password()
	}
	return
}
