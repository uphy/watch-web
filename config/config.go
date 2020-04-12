package config

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
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
	SourceConfig struct {
		DOM   *DOMSourceConfig   `json:"dom,omitempty"`
		Shell *ShellSourceConfig `json:"shell,omitempty"`
	}
	DOMSourceConfig struct {
		URL      TemplateString  `json:"url"`
		Selector TemplateString  `json:"selector"`
		Encoding *TemplateString `json:"encoding"`
	}
	ShellSourceConfig struct {
		Command *TemplateString `json:"command"`
	}
	ActionConfig struct {
		Slack      *SlackActionConfig      `json:"slack,omitempty"`
		LINENotify *LINENotifyActionConfig `json:"line_notify,omitempty"`
	}
	SlackActionConfig struct {
		URL TemplateString `json:"url"`
	}
	LINENotifyActionConfig struct {
		AccessToken TemplateString `json:"access_token"`
	}
	StoreConfig struct {
		Redis *RedisConfig `json:"redis,omitempty"`
	}
	RedisConfig struct {
		Address   *TemplateString `json:"address"`
		Password  *TemplateString `json:"password"`
		RedisToGo *TemplateString `json:"redistogo"`
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
