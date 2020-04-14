package config

import (
	"errors"
	"fmt"

	"github.com/uphy/watch-web/check"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
)

const (
	EmptyActionAccept EmptyAction = "accept"
	EmptyActionFail   EmptyAction = "fail"
)

type (
	EmptyAction  string
	SourceConfig struct {
		DOM         *DOMSourceConfig   `json:"dom,omitempty"`
		Shell       *ShellSourceConfig `json:"shell,omitempty"`
		EmptyAction *EmptyAction       `json:"empty,omitempty"`
		Retry       *int               `json:"retry,omitempty"`
	}
	DOMSourceConfig struct {
		URL      TemplateString  `json:"url"`
		Selector TemplateString  `json:"selector"`
		Encoding *TemplateString `json:"encoding"`
	}
	ShellSourceConfig struct {
		Command *TemplateString `json:"command"`
	}
	SourceWithRetry struct {
		source      check.Source
		emptyAction *EmptyAction
		retry       *int
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
	// wrap source for retry
	return &SourceWithRetry{source, s.EmptyAction, s.Retry}, nil
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

func (s *SourceWithRetry) Fetch() (string, error) {
	if s.retry == nil {
		return s.fetch()
	}
	retry := *s.retry
	var err error
	for i := 0; i <= retry; i++ {
		var v string
		v, err = s.fetch()
		if err != nil {
			continue
		}
		return v, nil
	}
	return "", fmt.Errorf("too many retries: lastError=%w", err)
}

func (s *SourceWithRetry) fetch() (string, error) {
	value, err := s.source.Fetch()
	if err != nil {
		return "", err
	}

	if s.emptyAction == nil || *s.emptyAction == EmptyActionAccept {
		return value, nil
	}
	if *s.emptyAction == EmptyActionFail {
		if len(value) == 0 {
			return "", errors.New("empty value")
		}
		return value, nil
	}
	return "", fmt.Errorf("unsupported empty action: %s", *s.emptyAction)
}
