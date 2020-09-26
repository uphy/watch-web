package config

import (
	"bytes"
	"io"

	"github.com/uphy/watch-web/pkg/domain/template"

	"github.com/ghodss/yaml"
)

type (
	Config struct {
		Jobs       []JobConfig              `json:"jobs"`
		InitialRun *template.TemplateString `json:"initial_run,omitempty"`
		Actions    []ActionConfig           `json:"actions"`
		Store      *StoreConfig             `json:"store"`
	}
)

func (c *Config) Save(w io.Writer) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, bytes.NewReader(data))
	return err
}
