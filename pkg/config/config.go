package config

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
	"github.com/uphy/watch-web/pkg/watch"
	"github.com/uphy/watch-web/pkg/config/template"
)

type (
	Config struct {
		Jobs       []JobConfig              `json:"jobs"`
		InitialRun *template.TemplateString `json:"initial_run,omitempty"`
		Actions    []ActionConfig           `json:"actions"`
		Store      *StoreConfig             `json:"store"`
	}
)

func LoadConfigFile(file string) (*Config, error) {
	// read file
	f, err := os.Open(file)
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

func (c *Config) Save(w io.Writer) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, bytes.NewReader(data))
	return err
}

func (c *Config) NewExecutor(log *logrus.Logger) (*watch.Executor, error) {
	ctx := template.NewRootContext()
	// store
	store, err := newStore(ctx, c.Store)
	if err != nil {
		return nil, err
	}
	log.WithFields(logrus.Fields{
		"store": fmt.Sprintf("%#v", store),
	}).Info("Created store.")

	// action
	actions := []watch.Action{}
	for _, actionConfig := range c.Actions {
		action, err := actionConfig.Action(ctx)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}

	// executor
	e := watch.NewExecutor(store, actions, log)
	if c.InitialRun != nil {
		initialRun, err := c.InitialRun.Evaluate(ctx)
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
		jobs, err := jobConfig.addTo(ctx, e)
		if err != nil {
			return nil, err
		}
		log.WithFields(logrus.Fields{
			"jobs": jobs,
		}).Debug("Added jobs to executor.")
	}
	return e, nil
}
