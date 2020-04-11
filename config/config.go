package config

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/uphy/watch-web/check"
)

type (
	Config struct {
		Jobs       map[string]JobConfig `json:"jobs"`
		InitialRun *bool                `json:"initial_run,omitempty"`
		Store      *StoreConfig         `json:"store"`
	}
	JobConfig struct {
		Label    string         `json:"label"`
		Link     string         `json:"link"`
		Source   *SourceConfig  `json:"source,omitempty"`
		Schedule string         `json:"schedule,omitempty"`
		Actions  []ActionConfig `json:"actions,omitempty"`
	}
	SourceConfig struct {
		DOM   *check.DOMSource   `json:"dom,omitempty"`
		Shell *check.ShellSource `json:"shell,omitempty"`
	}
	ActionConfig struct {
		Slack      *check.SlackAction      `json:"slack,omitempty"`
		LINENotify *check.LINENotifyAction `json:"line_notify,omitempty"`
	}
	StoreConfig struct {
		Redis *RedisConfig `json:"redis,omitempty"`
	}
	RedisConfig struct {
		Address   *string `json:"address"`
		Password  *string `json:"password"`
		RedisToGo *string `json:"redistogo"`
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

	// process template
	funcs := template.FuncMap{
		"env": func(name string) string {
			return os.Getenv(name)
		},
		"default": func(value string, defaultValue string) string {
			if value == "" {
				return defaultValue
			}
			return value
		},
		"sliceOf": func(values ...string) []string {
			return values
		},
	}
	tmpl, err := template.New("t").Funcs(funcs).Parse(string(data))
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, nil); err != nil {
		return nil, err
	}

	// parse yaml/json
	var config Config
	if err := yaml.Unmarshal(buf.Bytes(), &config); err != nil {
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

func (s *SourceConfig) Source() (check.Source, error) {
	if s.DOM != nil {
		return s.DOM, nil
	}
	if s.Shell != nil {
		return s.Shell, nil
	}
	return nil, errors.New("no source defined")
}

func (a *ActionConfig) Action() (check.Action, error) {
	if a.Slack != nil {
		return a.Slack, nil
	}
	if a.LINENotify != nil {
		return a.LINENotify, nil
	}
	return nil, errors.New("no action defined")
}
