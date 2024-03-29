package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/uphy/watch-web/pkg/domain/retry"
	"github.com/uphy/watch-web/pkg/domain/template"

	"github.com/uphy/watch-web/pkg/domain/value"

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
		log             *logrus.Logger
		ctx             *template.TemplateContext
		configDirectory *configDirectory
		store           domain.Store
	}
)

func LoadAndCreate(log *logrus.Logger, file string) (*watch.Executor, error) {
	baseDirectory, _ := filepath.Split(file)
	l := NewLoader(log, baseDirectory)
	conf, err := l.Load(file)
	if err != nil {
		return nil, err
	}
	return l.Create(conf)
}

func NewLoader(log *logrus.Logger, file string) *Loader {
	ctx := template.NewRootTemplateContext()
	dir, _ := filepath.Split(file)
	return &Loader{log, ctx, newConfigDirectory(dir), nil}
}

func (l *Loader) TemplateContext() *template.TemplateContext {
	return l.ctx
}

func (l *Loader) Load(file string) (*Config, error) {
	// read file
	resolved, err := l.configDirectory.resolve(file)
	if err != nil {
		return nil, fmt.Errorf("config file not found: file=%s", file)
	}
	f, err := os.Open(resolved)
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
	l.store = store
	l.log.WithFields(logrus.Fields{
		"store": fmt.Sprintf("%#v", store),
	}).Info("Created store.")

	// action
	actions, err := l.createActions(c.Actions)
	if err != nil {
		return nil, err
	}

	// executor
	e := watch.NewExecutor(store, l.log)
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
		if jobConfig.Enable != nil && !*jobConfig.Enable {
			continue
		}
		jobs, err := l.addJobTo(&jobConfig, e, actions)
		if err != nil {
			return nil, err
		}
		if len(jobs) == 0 {
			continue
		}
		l.log.WithFields(logrus.Fields{
			"jobs": jobs,
		}).Debug("Added jobs to executor.")
	}
	return e, nil
}

func (l *Loader) createStore(config *StoreConfig) (domain.Store, error) {
	if config == nil {
		return store.NewMemoryStore(), nil
	}
	if config.Redis != nil {
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
	if config.Directory != nil {
		return store.NewDirectoryStore(config.Directory.Path)
	}
	return nil, errors.New("unknown store")
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

func (l *Loader) createActions(a []ActionConfig) ([]domain.Action, error) {
	actions := []domain.Action{}
	for _, actionConfig := range a {
		action, err := l.createAction(&actionConfig)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}
	return actions, nil
}

func (l *Loader) createAction(a *ActionConfig) (domain.Action, error) {
	if a.SlackWebhook != nil {
		return l.createActionSlackWebhook(a.SlackWebhook)
	}
	if a.SlackBot != nil {
		return l.createActionSlackBot(a.SlackBot)
	}
	if a.Console != nil {
		return action.NewConsoleAction(), nil
	}
	return nil, errors.New("no action defined")
}

func (l *Loader) createActionSlackWebhook(s *SlackWebhookActionConfig) (domain.Action, error) {
	url, err := s.URL.Evaluate(l.ctx)
	if err != nil {
		return nil, err
	}
	return action.NewSlackWebhookAction(url, s.Debug), nil
}

func (l *Loader) createActionSlackBot(s *SlackBotActionConfig) (domain.Action, error) {
	channel, err := s.Channel.Evaluate(l.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate slack channel template: %w", err)
	}
	token, err := s.Token.Evaluate(l.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate slack token template: %w", err)
	}
	repo, err := s.newRepository(l.ctx, l.store)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository for slack bot action: %w", err)
	}
	return action.NewSlackBotAction(token, channel, s.Debug, repo), nil
}

func (l *Loader) addJobTo(c *JobConfig, e *watch.Executor, actions []domain.Action) ([]*watch.Job, error) {
	jobs := make([]*watch.Job, 0)
	if c.WithItems == nil {
		job, err := l.addJobOne(c, e, actions)
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
			j, err := l.addJobOne(c, e, actions)
			if err != nil {
				return nil, err
			}
			jobs = append(jobs, j)
			l.ctx.PopScope()
		}
	}
	return jobs, nil
}

func evaluateItemAsTemplate(ctx *template.TemplateContext, v interface{}) (interface{}, error) {
	m, ok := v.(map[string]interface{})
	if ok {
		evaluated := make(map[string]interface{})
		for key, value := range m {
			ekey, err := template.TemplateString(key).Evaluate(ctx)
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
	e, err := template.TemplateString(fmt.Sprint(v)).Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (l *Loader) addJobOne(c *JobConfig, e *watch.Executor, actions []domain.Action) (*watch.Job, error) {
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
	jobActions, err := l.createActions(c.Actions)
	if err != nil {
		return nil, err
	}
	if len(jobActions) > 0 {
		actions = jobActions
	}
	job := watch.NewJob(&domain.JobInfo{
		ID:    id,
		Label: label,
		Link:  link,
	}, source, actions)

	if err := e.AddJob(job, &schedule); err != nil {
		return nil, err
	}
	return job, nil
}

func (l *Loader) CreateSource(s *SourceConfig) (domain.Source, error) {
	// load raw source
	var src domain.Source
	var err error
	// constant source must be evaluated first for include source's override config.
	if s.Constant != nil {
		src, err = l.createSourceConstant(s.Constant)
	} else if s.DOM != nil {
		src, err = l.createSourceDOM(s.DOM)
	} else if s.Shell != nil {
		src, err = l.createSourceShell(s.Shell)
	} else if s.Include != nil {
		src, err = l.createSourceInclude(s.Include)
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
	var retrier *retry.Retrier
	if s.Retry != nil {
		r, err := parseRetry(s.Retry)
		if err != nil {
			return nil, err
		}
		retrier = r
	}
	return source.NewRetrySource(src, s.EmptyAction, retrier), nil
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
		v, err := value.ConvertInterfaceAs(s.Value, value.ValueTypeAutoDetect)
		if err != nil {
			return nil, err
		}
		return source.NewConstantSource(v), nil
	}
	if s.File != nil {
		dir := l.configDirectory.childRelative("constants")
		file, err := dir.resolve(*s.File)
		if err != nil {
			return nil, fmt.Errorf("cannot resolve constant file: %w", err)
		}
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		text, err := template.TemplateString(string(b)).Evaluate(l.ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate template of constant source file: %w", err)
		}
		return source.NewConstantSource(value.NewStringValue(text)), nil
	}
	if s.Template != nil {
		v, err := s.Template.Evaluate(l.ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate constant source template: %w", err)
		}
		return source.NewConstantSource(value.NewStringValue(v)), nil
	}
	return nil, errors.New("unsupported constant source")
}

func (l *Loader) createSourceInclude(c *IncludeSourceConfig) (domain.Source, error) {
	// resolve file
	filePath, err := c.File.Evaluate(l.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate include file path template: file=%s, err=%w", filePath, err)
	}
	dir := l.configDirectory.childRelative("includes")
	resolvedFilePath, err := dir.resolve(filePath)
	if err != nil {
		return nil, fmt.Errorf("include file not found: file=%s", filePath)
	}
	// open
	f, err := os.Open(resolvedFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open include file: file=%s, err=%w", resolvedFilePath, err)
	}
	defer f.Close()
	// read as yaml
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var srcConfig *SourceConfig
	if err := yaml.Unmarshal(b, &srcConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal include yaml file: file=%s, err=%w", resolvedFilePath, err)
	}
	// initialize template variables for override config and source
	l.ctx.PushScope()
	defer l.ctx.PopScope()
	for k, v := range c.Vars {
		evaluated, err := v.Evaluate(l.ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate include var template: key=%s, template=%v, err=%w", k, v, err)
		}
		l.ctx.Set(k, evaluated)
	}
	// override config
	if c.Overrides != nil && c.Overrides.Constant != nil {
		srcConfig.Constant = c.Overrides.Constant
		// srcConfig may have another source such as dom, shell
		// but I don't need to clear them because they won't be evaluated
		// because constant source is evaluated first in CreateSource().
	}
	// create source
	src, err := l.CreateSource(srcConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create included source: file=%s, err=%w", resolvedFilePath, err)
	}
	return src, nil
}

func (l *Loader) createTransforms(t TransformsConfig, src domain.Source) (domain.Source, error) {
	if len(t) == 0 {
		return src, nil
	}
	transformers := make([]domain.Transformer, 0)
	for _, transformConfig := range t {
		t, err := l.createTransform(&transformConfig)
		if err != nil {
			return nil, err
		}
		if transformConfig.Retry != nil {
			retrier, err := parseRetry(transformConfig.Retry)
			if err != nil {
				return nil, err
			}
			t = transformer.NewRetryTransform(t, retrier)
		}
		transformers = append(transformers, t)
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
	if t.Map != nil {
		return transformer.NewMapTransformer(t.Map.Template, l.ctx.Snapshot()), nil
	}
	if t.JSONObject != nil {
		return transformer.NewJSONObjectTransformer(), nil
	}
	if t.JSONArray != nil {
		return transformer.NewJSONArrayTransformer(l.ctx.Snapshot(), t.JSONArray.Condition), nil
	}
	if t.Sort != nil {
		return transformer.NewSortTransformer(t.Sort.By), nil
	}
	if t.Script != nil {
		scr, err := t.Script.NewScript(l.ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to parse script: %v", err)
		}
		return transformer.NewScriptTransformer(scr)
	}
	if t.Filter != nil {
		scr, err := t.Filter.NewScript(l.ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to parse script: %v", err)
		}
		return transformer.NewFilterTransformer(scr), nil
	}
	if t.Debug != nil {
		return transformer.NewDebugTransformer(*t.Debug), nil
	}
	return nil, errors.New("no transforms defined")
}
