package config

import (
	"errors"
	"fmt"

	"github.com/uphy/watch-web/check"
)

type (
	FiltersConfig []FilterConfig
	FilterConfig  struct {
		Template *TemplateString `json:"template,omitempty"`
	}

	FilterSource struct {
		source  check.Source
		filters []Filter
	}
	Filter interface {
		Filter(s string) (string, error)
	}
	TemplateFilter struct {
		template TemplateString
		ctx      *TemplateContext
	}
)

func (f FiltersConfig) Filters(ctx *TemplateContext, source check.Source) (check.Source, error) {
	if len(f) == 0 {
		return source, nil
	}
	filters := make([]Filter, 0)
	for _, filterConfig := range f {
		filter, err := filterConfig.Filter(ctx)
		if err != nil {
			return nil, err
		}
		filters = append(filters, filter)
	}
	return &FilterSource{source, filters}, nil
}

func (f *FilterConfig) Filter(ctx *TemplateContext) (Filter, error) {
	if f.Template != nil {
		return &TemplateFilter{*f.Template, ctx}, nil
	}
	return nil, errors.New("no filters defined")
}

func (f *FilterSource) Fetch() (string, error) {
	s, err := f.source.Fetch()
	if err != nil {
		return "", err
	}
	for _, filter := range f.filters {
		filtered, err := filter.Filter(s)
		if err != nil {
			return "", err
		}
		s = filtered
	}
	return s, nil
}

func (f *FilterSource) String() string {
	return fmt.Sprintf("Filter[source=%v, filters=%v]", f.source, f.filters)
}

func (t *TemplateFilter) Filter(s string) (string, error) {
	t.ctx.PushScope()
	defer t.ctx.PopScope()
	t.ctx.Set("source", s)
	evaluated, err := t.template.Evaluate(t.ctx)
	if err != nil {
		return "", err
	}
	return evaluated, nil
}

func (t *TemplateFilter) String() string {
	return fmt.Sprintf("Template[template=%v]", t.template)
}
