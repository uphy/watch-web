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
	"github.com/uphy/watch-web/pkg/check"
)

type (
	Config struct {
		Jobs       []JobConfig     `json:"jobs"`
		InitialRun *TemplateString `json:"initial_run,omitempty"`
		Store      *StoreConfig    `json:"store"`
	}
	JobConfig struct {
		ID        TemplateString `json:"id"`
		Label     TemplateString `json:"label"`
		Link      TemplateString `json:"link"`
		Source    *SourceConfig  `json:"source,omitempty"`
		Schedule  TemplateString `json:"schedule,omitempty"`
		Actions   []ActionConfig `json:"actions"`
		WithItems []interface{}  `json:"with_items,omitempty"`
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

func (c *Config) NewExecutor(log *logrus.Logger) (*check.Executor, error) {
	ctx := NewRootContext()
	store, err := newStore(ctx, c.Store)
	if err != nil {
		return nil, err
	}
	log.WithFields(logrus.Fields{
		"store": fmt.Sprintf("%#v", store),
	}).Info("Created store.")
	e := check.NewExecutor(store, log)
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
