package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/uphy/watch-web/pkg/domain"
	"github.com/uphy/watch-web/pkg/watch/source"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
)

type (
	SourceConfig struct {
		DOM        *DOMSourceConfig      `json:"dom,omitempty"`
		Shell      *ShellSourceConfig    `json:"shell,omitempty"`
		Constant   *ConstantSourceConfig `json:"constant",omitempty`
		Transforms TransformsConfig      `json:"transforms,omitempty"`

		EmptyAction *source.EmptyAction `json:"empty,omitempty"`
		Retry       *int                `json:"retry,omitempty"`
	}
	DOMSourceConfig struct {
		URL      domain.TemplateString  `json:"url"`
		Selector domain.TemplateString  `json:"selector"`
		Encoding *domain.TemplateString `json:"encoding"`
	}
	ShellSourceConfig struct {
		Command *domain.TemplateString `json:"command"`
	}
	ConstantSourceConfig struct {
		Value    interface{}            `json:"value,omitempty"`
		Template *domain.TemplateString `json:"template,omitempty"`
		File     *string                `json:"file,omitempty"`
	}
)

func (s *SourceConfig) Source(ctx *domain.TemplateContext) (domain.Source, error) {
	// load raw source
	var src domain.Source
	var err error
	if s.DOM != nil {
		src, err = s.DOM.Source(ctx)
	} else if s.Shell != nil {
		src, err = s.Shell.Source(ctx)
	} else if s.Constant != nil {
		src, err = s.Constant.Source(ctx)
	}
	if err != nil {
		return nil, err
	}
	if src == nil {
		return nil, errors.New("no source defined")
	}
	// wrap source for transformers
	if len(s.Transforms) > 0 {
		src, err = s.Transforms.Transforms(ctx, src)
		if err != nil {
			return nil, err
		}
	}
	// wrap source for retry
	return source.NewRetrySource(src, s.EmptyAction, s.Retry), nil
}

func (d *DOMSourceConfig) Source(ctx *domain.TemplateContext) (domain.Source, error) {
	var encoding encoding.Encoding
	if d.Encoding != nil {
		enc, err := d.Encoding.Evaluate(ctx)
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
	u, err := d.URL.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	s, err := d.Selector.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	source := source.NewDOMSource(u, s, encoding)
	return source, nil
}

func (d *ShellSourceConfig) Source(ctx *domain.TemplateContext) (domain.Source, error) {
	command, err := d.Command.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return source.NewShellSource(command), nil
}

func (s *ConstantSourceConfig) Source(ctx *domain.TemplateContext) (domain.Source, error) {
	if s.Value != nil {
		v, err := domain.ConvertInterfaceAs(s.Value, domain.ValueTypeAutoDetect)
		if err != nil {
			return nil, err
		}
		return source.NewConstantSource(v), nil
	}
	if s.File != nil {
		f, err := os.Open(*s.File)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		return source.NewConstantSource(domain.NewStringValue(string(b))), nil
	}
	if s.Template != nil {
		value, err := s.Template.Evaluate(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate constant source template: %w", err)
		}
		return source.NewConstantSource(domain.NewStringValue(string(value))), nil
	}
	return nil, errors.New("unsupported constant source")
}
