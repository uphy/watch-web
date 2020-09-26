package config

import (
	"fmt"

	"github.com/uphy/watch-web/pkg/domain/retry"
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
		Retry       interface{}         `json:"retry,omitempty"`
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
	Retry struct {
		Retry               int      `json:"retry"`
		InitialInterval     *float64 `json:"initial_interval,omitempty"`
		Multiplier          *float64 `json:"multiplier,omitempty"`
		RandomizationFactor *float64 `json:"randomization_factor,omitempty"`
		MaxInterval         *float64 `json:"max_interval"`
	}
)

func parseRetry(r interface{}) (*retry.Retrier, error) {
	switch re := r.(type) {
	case *Retry:
		b := retry.NewBuilder(re.Retry)
		if re.InitialInterval != nil {
			b.InitialInterval(*re.InitialInterval)
		}
		if re.MaxInterval != nil {
			b.MaxInterval(*re.MaxInterval)
		}
		if re.RandomizationFactor != nil {
			b.RandomizationFactor(*re.RandomizationFactor)
		}
		if re.Multiplier != nil {
			b.Multiplier(*re.Multiplier)
		}
		return b.Build(), nil
	case float64:
		return retry.NewBuilder(int(re)).Build(), nil
	}
	return nil, fmt.Errorf("unsupported retry: %v", r)
}
