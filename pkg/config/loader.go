package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"

	"github.com/ghodss/yaml"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"github.com/uphy/watch-web/pkg/domain"
	"github.com/uphy/watch-web/pkg/watch"
	"github.com/uphy/watch-web/pkg/watch/action"
	"github.com/uphy/watch-web/pkg/watch/source"
	"github.com/uphy/watch-web/pkg/watch/store"
	"github.com/uphy/watch-web/pkg/watch/transformer"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
)

type (
	Loader struct {
		log  *logrus.Logger
		file string
		ctx  *domain.TemplateContext
	}
)

func LoadAndCreate(log *logrus.Logger, file string) (*watch.Executor, error) {
	l := NewLoader(log, file)
	conf, err := l.Load()
	if err != nil {
		return nil, err
	}
	return l.Create(conf)
}

func NewLoader(log *logrus.Logger, file string) *Loader {
	ctx := domain.NewRootTemplateContext()
	return &Loader{log, file, ctx}
}

func (l *Loader) TemplateContext() *domain.TemplateContext {
	return l.ctx
}

func (l *Loader) Load() (*Config, error) {
	// read file
	f, err := os.Open(l.file)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	// parse yaml/json
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (l *Loader) Create(c *Config) (*watch.Executor, error) {
	// store
	store, err := l.createStore(c.Store)
	if err != nil {
		return nil, err
	}
	l.log.WithFields(logrus.Fields{
		"store": fmt.Sprintf("%#v", store),
	}).Info("Created store.")

	// action
	actions := []domain.Action{}
	for _, actionConfig := range c.Actions {
		action, err := l.createAction(&actionConfig)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}

	// executor
	e := watch.NewExecutor(store, actions, l.log)
	if c.InitialRun != nil {
		initialRun, err := c.InitialRun.Evaluate(l.ctx)
		if err != nil {
			return nil, err
		}
		ini, err := strconv.ParseBool(initialRun)
		if err != nil {
			return nil, err
		}
		e.InitialRun = ini
	}

	// jobs
	for _, jobConfig := range c.Jobs {
		jobs, err := l.addJobTo(&jobConfig, e)
		if err != nil {
			return nil, err
		}
		l.log.WithFields(logrus.Fields{
			"jobs": jobs,
		}).Debug("Added jobs to executor.")
	}
	return e, nil
}

func (l *Loader) createStore(config *StoreConfig) (domain.Store, error) {
	if config != nil && config.Redis != nil {
		password := ""
		addr := ""
		if config.Redis.Address != nil {
			if config.Redis.Password != nil {
				p, err := config.Redis.Password.Evaluate(l.ctx)
				if err != nil {
					return nil, err
				}
				password = p
			}
			a, err := config.Redis.Address.Evaluate(l.ctx)
			if err != nil {
				return nil, err
			}
			addr = a
		} else if config.Redis.RedisToGo != nil {
			r, err := config.Redis.RedisToGo.Evaluate(l.ctx)
			if err != nil {
				return nil, err
			}
			addr, password, err = parseRedisToGoURL(r)
			if err != nil {
				return nil, err
			}
		}
		if addr != "" {
			client := redis.NewClient(&redis.Options{
				Addr:     addr,
				Password: password,
			})
			return store.NewRedisStore(client), nil
		}
	}
	return store.NewMemoryStore(), nil
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

func (l *Loader) createAction(a *ActionConfig) (domain.Action, error) {
	if a.Slack != nil {
		return l.createActionSlack(a.Slack)
	}
	if a.Console != nil {
		return action.NewConsoleAction(), nil
	}
	return nil, errors.New("no action defined")
}

func (l *Loader) createActionSlack(s *SlackActionConfig) (domain.Action, error) {
	url, err := s.URL.Evaluate(l.ctx)
	if err != nil {
		return nil, err
	}
	return action.NewSlackAction(url, s.Debug), nil
}

func (l *Loader) addJobTo(c *JobConfig, e *watch.Executor) ([]*watch.Job, error) {
	jobs := make([]*watch.Job, 0)
	if len(c.WithItems) == 0 {
		job, err := l.addJobOne(c, e)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	} else {
		for itemIndex, item := range c.WithItems {
			evaluatedItem, err := evaluateItemAsTemplate(l.ctx, item)
			if err != nil {
				return nil, err
			}
			l.ctx.PushScope()
			l.ctx.Set("itemIndex", itemIndex)
			l.ctx.Set("item", evaluatedItem)
			j, err := l.addJobOne(c, e)
			if err != nil {
				return nil, err
			}
			jobs = append(jobs, j)
			l.ctx.PopScope()
		}
	}
	return jobs, nil
}

func evaluateItemAsTemplate(ctx *domain.TemplateContext, v interface{}) (interface{}, error) {
	m, ok := v.(map[string]interface{})
	if ok {
		evaluated := make(map[string]interface{})
		for key, value := range m {
			ekey, err := domain.TemplateString(key).Evaluate(ctx)
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
	e, err := domain.TemplateString(fmt.Sprint(v)).Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (l *Loader) addJobOne(c *JobConfig, e *watch.Executor) (*watch.Job, error) {
	source, err := l.CreateSource(c.Source)
	if err != nil {
		return nil, err
	}
	id, err := c.ID.Evaluate(l.ctx)
	if err != nil {
		return nil, err
	}
	schedule, err := c.Schedule.Evaluate(l.ctx)
	if err != nil {
		return nil, err
	}
	label, err := c.Label.Evaluate(l.ctx)
	if err != nil {
		return nil, err
	}
	link, err := c.Link.Evaluate(l.ctx)
	if err != nil {
		return nil, err
	}
	job := watch.NewJob(&domain.JobInfo{
		ID:    id,
		Label: label,
		Link:  link,
	}, source)

	if err := e.AddJob(job, &schedule); err != nil {
		return nil, err
	}
	return job, nil
}

func (l *Loader) CreateSource(s *SourceConfig) (domain.Source, error) {
	// load raw source
	var src domain.Source
	var err error
	if s.DOM != nil {
		src, err = l.createSourceDOM(s.DOM)
	} else if s.Shell != nil {
		src, err = l.createSourceShell(s.Shell)
	} else if s.Constant != nil {
		src, err = l.createSourceConstant(s.Constant)
	}
	if err != nil {
		return nil, err
	}
	if src == nil {
		return nil, errors.New("no source defined")
	}
	// wrap source for transformers
	if len(s.Transforms) > 0 {
		src, err = l.createTransforms(s.Transforms, src)
		if err != nil {
			return nil, err
		}
	}
	// wrap source for retry
	return source.NewRetrySource(src, s.EmptyAction, s.Retry), nil
}

func (l *Loader) createSourceDOM(d *DOMSourceConfig) (domain.Source, error) {
	var encoding encoding.Encoding
	if d.Encoding != nil {
		enc, err := d.Encoding.Evaluate(l.ctx)
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
	u, err := d.URL.Evaluate(l.ctx)
	if err != nil {
		return nil, err
	}
	s, err := d.Selector.Evaluate(l.ctx)
	if err != nil {
		return nil, err
	}
	source := source.NewDOMSource(u, s, encoding)
	return source, nil
}

func (l *Loader) createSourceShell(d *ShellSourceConfig) (domain.Source, error) {
	command, err := d.Command.Evaluate(l.ctx)
	if err != nil {
		return nil, err
	}
	return source.NewShellSource(command), nil
}

func (l *Loader) createSourceConstant(s *ConstantSourceConfig) (domain.Source, error) {
	if s.Value != nil {
		v, err := domain.ConvertInterfaceAs(s.Value, domain.ValueTypeAutoDetect)
		if err != nil {
			return nil, err
		}
		return source.NewConstantSource(v), nil
	}
	if s.File != nil {
		f, err := os.Open(*s.File)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		return source.NewConstantSource(domain.NewStringValue(string(b))), nil
	}
	if s.Template != nil {
		value, err := s.Template.Evaluate(l.ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate constant source template: %w", err)
		}
		return source.NewConstantSource(domain.NewStringValue(string(value))), nil
	}
	return nil, errors.New("unsupported constant source")
}

func (l *Loader) createTransforms(t TransformsConfig, src domain.Source) (domain.Source, error) {
	if len(t) == 0 {
		return src, nil
	}
	transformers := make([]domain.Transformer, 0)
	for _, transformConfig := range t {
		transformer, err := l.createTransform(&transformConfig)
		if err != nil {
			return nil, err
		}
		transformers = append(transformers, transformer)
	}
	return source.NewTransformerSource(src, transformers), nil
}

func (l *Loader) createTransform(t *TransformConfig) (domain.Transformer, error) {
	if t.Template != nil {
		return transformer.NewTemplateTransformer(*t.Template, l.ctx), nil
	}
	if t.DOM != nil {
		selector, err := t.DOM.Evaluate(l.ctx.Snapshot())
		if err != nil {
			return nil, err
		}
		return transformer.NewDOMTransformer(selector), nil
	}
	if t.JSONTransform != nil {
		return transformer.NewJSONTransformTransformer(t.JSONTransform.Transform, l.ctx.Snapshot()), nil
	}
	if t.JSONObject != nil {
		return transformer.NewJSONObjectTransformer(), nil
	}
	if t.JSONArray != nil {
		return transformer.NewJSONArrayTransformer(l.ctx.Snapshot(), t.JSONArray.Condition), nil
	}
	if t.JSONArraySort != nil {
		return transformer.NewJSONArraySortTransformer(t.JSONArraySort.By), nil
	}
	if t.Script != nil {
		script, err := t.Script.Script.Evaluate(l.ctx)
		if err != nil {
			return nil, err
		}
		language := "anko"
		if t.Script.Language != nil {
			language = *t.Script.Language
		}
		switch language {
		case "anko":
			return transformer.NewScriptTransformer(script)
		default:
			return nil, fmt.Errorf("unsupported script language: %s", language)
		}
	}
	if t.Debug != nil {
		return transformer.NewDebugTransformer(*t.Debug), nil
	}
	return nil, errors.New("no transforms defined")
}
