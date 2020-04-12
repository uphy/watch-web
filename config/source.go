package config

import (
	"errors"

	"github.com/uphy/watch-web/check"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
)

type (
	SourceConfig struct {
		DOM      *DOMSourceConfig   `json:"dom,omitempty"`
		Shell    *ShellSourceConfig `json:"shell,omitempty"`
		Template *TemplateString    `json:"template,omitempty"`
	}
	DOMSourceConfig struct {
		URL      TemplateString  `json:"url"`
		Selector TemplateString  `json:"selector"`
		Encoding *TemplateString `json:"encoding"`
	}
	ShellSourceConfig struct {
		Command *TemplateString `json:"command"`
	}
	TemplateSource struct {
		template TemplateString
		ctx      *TemplateContext
		source   check.Source
	}
)

func (s *SourceConfig) Source(ctx *TemplateContext) (check.Source, error) {
	// load raw source
	var source check.Source
	var err error
	if s.DOM != nil {
		source, err = s.DOM.Source(ctx)
	} else if s.Shell != nil {
		source, err = s.Shell.Source(ctx)
	}
	if err != nil {
		return nil, err
	}
	if source == nil {
		return nil, errors.New("no source defined")
	}

	// wrap with template source
	if s.Template != nil {
		source = &TemplateSource{*s.Template, ctx, source}
	}
	return source, nil
}

func (d *DOMSourceConfig) Source(ctx *TemplateContext) (check.Source, error) {
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
	source := check.NewDOMSource(u, s, encoding)
	return source, nil
}

func (d *ShellSourceConfig) Source(ctx *TemplateContext) (check.Source, error) {
	command, err := d.Command.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return check.NewShellSource(command), nil
}

func (t *TemplateSource) Fetch() (string, error) {
	s, err := t.source.Fetch()
	if err != nil {
		return "", err
	}
	t.ctx.PushScope()
	defer t.ctx.PopScope()
	t.ctx.Set("output", s)
	return t.template.Evaluate(t.ctx)
}
