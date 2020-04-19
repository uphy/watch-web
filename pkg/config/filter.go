package config

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/uphy/watch-web/pkg/watch"
	"github.com/uphy/watch-web/pkg/config/filter"
	"github.com/uphy/watch-web/pkg/config/template"
	"github.com/uphy/watch-web/pkg/value"
)

type (
	FiltersConfig []FilterConfig
	FilterConfig  struct {
		Template      *template.TemplateString `json:"template,omitempty"`
		DOM           *template.TemplateString `json:"dom,omitempty"`
		JSONArray     *struct{}                `json:"json_array,omitempty"`
		JSONObject    *struct{}                `json:"json_object,omitempty"`
		JSONTransform *struct {
			Transform map[string]template.TemplateString `json:"transform,omitempty"`
		} `json:"json_transform,omitempty"`
		JSONArraySort *struct {
			By string `json:"by"`
		} `json:"json_array_sort"`
	}

	FilterSource struct {
		source  watch.Source
		filters []Filter
	}
	Filter interface {
		Filter(ctx *watch.JobContext, v value.Value) (value.Value, error)
	}
)

func (f FiltersConfig) Filters(ctx *template.TemplateContext, source watch.Source) (watch.Source, error) {
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

func (f *FilterConfig) Filter(ctx *template.TemplateContext) (Filter, error) {
	if f.Template != nil {
		return filter.NewTemplateFilter(*f.Template, ctx), nil
	}
	if f.DOM != nil {
		selector, err := f.DOM.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		return filter.NewDOMFilter(selector), nil
	}
	if f.JSONTransform != nil {
		return filter.NewJSONTransformFilter(f.JSONTransform.Transform, ctx), nil
	}
	if f.JSONObject != nil {
		return filter.NewJSONObjectFilter(), nil
	}
	if f.JSONArray != nil {
		return filter.NewJSONArrayFilter(), nil
	}
	if f.JSONArraySort != nil {
		return filter.NewJSONArraySortFilter(f.JSONArraySort.By), nil
	}
	return nil, errors.New("no filters defined")
}

func (f *FilterSource) Fetch(ctx *watch.JobContext) (value.Value, error) {
	v, err := f.source.Fetch(ctx)
	if err != nil {
		return nil, err
	}
	ctx.Log.WithField("source", v).Debug("Start filter chain.")
	for _, filter := range f.filters {
		filtered, err := filter.Filter(ctx, v)
		if err != nil {
			ctx.Log.WithFields(logrus.Fields{
				"filter": fmt.Sprintf("%#v", filter),
			}).Debug("Failed to filtered value.")
			return nil, err
		}
		ctx.Log.WithFields(logrus.Fields{
			"filter":   fmt.Sprintf("%#v", filter),
			"filtered": filtered,
		}).Debug("Filtered value.")
		v = filtered
	}
	return v, nil
}

func (f *FilterSource) String() string {
	return fmt.Sprintf("Filter[source=%v, filters=%v]", f.source, f.filters)
}
