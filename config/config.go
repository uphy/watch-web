package config

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/uphy/watch-web/check"
)

type (
	Config struct {
		Jobs       []JobConfig  `json:"jobs"`
		InitialRun *bool        `json:"initial_run,omitempty"`
		Store      *StoreConfig `json:"store"`
	}
	JobConfig struct {
		ID        TemplateString `json:"id"`
		Label     TemplateString `json:"label"`
		Link      TemplateString `json:"link"`
		Source    *SourceConfig  `json:"source,omitempty"`
		Schedule  TemplateString `json:"schedule,omitempty"`
		Actions   []ActionConfig `json:"actions,omitempty"`
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
		jobConfig.addTo(ctx, e)
	}
	return e, nil
}
