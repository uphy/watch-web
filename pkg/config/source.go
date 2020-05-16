package config

import (
	"github.com/uphy/watch-web/pkg/domain/template"
	"github.com/uphy/watch-web/pkg/watch/source"
)

type (
	SourceConfig struct {
		DOM        *DOMSourceConfig      `json:"dom,omitempty"`
		Shell      *ShellSourceConfig    `json:"shell,omitempty"`
		Constant   *ConstantSourceConfig `json:"constant,omitempty"`
		Include    *IncludeSourceConfig  `json:"include,omitempty"`
		Transforms TransformsConfig      `json:"transforms,omitempty"`

		EmptyAction *source.EmptyAction `json:"empty,omitempty"`
		Retry       *int                `json:"retry,omitempty"`
	}
	DOMSourceConfig struct {
		URL      template.TemplateString  `json:"url"`
		Selector template.TemplateString  `json:"selector"`
		Encoding *template.TemplateString `json:"encoding"`
	}
	ShellSourceConfig struct {
		Command *template.TemplateString `json:"command"`
	}
	ConstantSourceConfig struct {
		Value    interface{}              `json:"value,omitempty"`
		Template *template.TemplateString `json:"template,omitempty"`
		File     *string                  `json:"file,omitempty"`
	}
	IncludeSourceConfig struct {
		File template.TemplateString `json:"file"`
		// Overrides defines a source config.
		// `File`'s source will be overriden by this source.
		// This is for testing
		Overrides *SourceConfig                      `json:"overrides,omitempty"`
		Vars      map[string]template.TemplateString `json:"vars,omitempty"`
	}
)
