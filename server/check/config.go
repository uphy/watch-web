package check

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/ghodss/yaml"
)

type (
	Config struct {
		Jobs       map[string]JobConfig `json:"jobs"`
		InitialRun *bool                `json:"initial_run,omitempty"`
	}
	JobConfig struct {
		Source   *SourceConfig  `json:"source,omitempty"`
		Schedule string         `json:"schedule,omitempty"`
		Actions  []ActionConfig `json:"actions,omitempty"`
	}
	SourceConfig struct {
		DOM     *DOMSource     `json:"dom,omitempty"`
		Command *CommandSource `json:"command,omitempty"`
	}
	ActionConfig struct {
		Slack      *SlackAction      `json:"slack,omitempty"`
		LINENotify *LINENotifyAction `json:"line_notify,omitempty"`
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

func (c *Config) NewExecutor() (*Executor, error) {
	e := NewExecutor()
	if c.InitialRun != nil {
		e.initialRun = *c.InitialRun
	}
	for name, job := range c.Jobs {
		source, err := job.Source.Source()
		if err != nil {
			return nil, err
		}
		actions := []Action{}
		for _, a := range job.Actions {
			action, err := a.Action()
			if err != nil {
				return nil, err
			}
			actions = append(actions, action)
		}
		e.AddJob(name, job.Schedule, source, actions...)
	}
	return e, nil
}

func (s *SourceConfig) Source() (Source, error) {
	if s.DOM != nil {
		return s.DOM, nil
	}
	if s.Command != nil {
		return s.Command, nil
	}
	return nil, errors.New("no source defined")
}

func (a *ActionConfig) Action() (Action, error) {
	if a.Slack != nil {
		return a.Slack, nil
	}
	if a.LINENotify != nil {
		return a.LINENotify, nil
	}
	return nil, errors.New("no action defined")
}
