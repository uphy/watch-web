package config

import (
	"errors"

	"github.com/uphy/watch-web/check"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
)

type (
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
)

func (s *SourceConfig) Source(ctx *TemplateContext) (check.Source, error) {
	if s.DOM != nil {
		return s.DOM.Source(ctx)
	}
	if s.Shell != nil {
		return s.Shell.Source(ctx)
	}
	return nil, errors.New("no source defined")
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
