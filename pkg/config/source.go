package config

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/uphy/watch-web/pkg/check"
	"github.com/uphy/watch-web/pkg/value"
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
		DOM     *DOMSourceConfig   `json:"dom,omitempty"`
		Shell   *ShellSourceConfig `json:"shell,omitempty"`
		Filters FiltersConfig      `json:"filters,omitempty"`

		EmptyAction *EmptyAction `json:"empty,omitempty"`
		Retry       *int         `json:"retry,omitempty"`
	}
	DOMSourceConfig struct {
		URL      TemplateString  `json:"url"`
		Selector TemplateString  `json:"selector"`
		Encoding *TemplateString `json:"encoding"`
	}
	ShellSourceConfig struct {
		Command *TemplateString  `json:"command"`
		Type    *value.ValueType `json:"type,omitempty"`
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
	// wrap source for filters
	if len(s.Filters) > 0 {
		source, err = s.Filters.Filters(ctx, source)
		if err != nil {
			return nil, err
		}
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
	return check.NewShellSource(command, d.valueType()), nil
}

func (s *ShellSourceConfig) valueType() value.ValueType {
	if s.Type == nil {
		return value.ValueTypeString
	}
	return *s.Type
}

func (s *SourceWithRetry) Fetch(ctx *check.JobContext) (value.Value, error) {
	if s.retry == nil {
		return s.fetch(ctx)
	}
	retry := *s.retry
	var err error
	for i := 0; i <= retry; i++ {
		var v value.Value
		v, err = s.fetch(ctx)
		if err == nil {
			return v, nil
		}

		if i < retry {
			waitSeconds := 1 + int(math.Pow(2, float64(i)))
			time.Sleep(time.Second * time.Duration(waitSeconds))
		}
	}
	return nil, fmt.Errorf("too many retries: lastError=%w", err)
}

func (s *SourceWithRetry) fetch(ctx *check.JobContext) (value.Value, error) {
	v, err := s.source.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	if s.emptyAction == nil || *s.emptyAction == EmptyActionAccept {
		return v, nil
	}
	if *s.emptyAction == EmptyActionFail {
		if v.Empty() {
			return nil, errors.New("empty value")
		}
		return v, nil
	}
	return nil, fmt.Errorf("unsupported empty action: %s", *s.emptyAction)
}

func (s *SourceWithRetry) String() string {
	var emptyAction = ""
	var retry = ""
	if s.emptyAction != nil {
		emptyAction = string(*s.emptyAction)
	}
	if s.retry != nil {
		retry = fmt.Sprint(*s.retry)
	}
	return fmt.Sprintf("Retry[source=%v, retry=%v, emptyAction=%v]", s.source, retry, emptyAction)
}
