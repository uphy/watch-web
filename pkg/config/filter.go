package config

import (
	"errors"

	"github.com/uphy/watch-web/pkg/domain"
	"github.com/uphy/watch-web/pkg/watch/filter"
	"github.com/uphy/watch-web/pkg/watch/source"
)

type (
	FiltersConfig []FilterConfig
	FilterConfig  struct {
		Template  *domain.TemplateString `json:"template,omitempty"`
		DOM       *domain.TemplateString `json:"dom,omitempty"`
		JSONArray *struct {
			Condition *domain.TemplateString `json:"condition,omitempty"`
		} `json:"json_array,omitempty"`
		JSONObject    *struct{} `json:"json_object,omitempty"`
		JSONTransform *struct {
			Transform map[string]domain.TemplateString `json:"transform,omitempty"`
		} `json:"json_transform,omitempty"`
		JSONArraySort *struct {
			By string `json:"by"`
		} `json:"json_array_sort"`
	}
)

func (f FiltersConfig) Filters(ctx *domain.TemplateContext, src domain.Source) (domain.Source, error) {
	if len(f) == 0 {
		return src, nil
	}
	filters := make([]domain.Filter, 0)
	for _, filterConfig := range f {
		filter, err := filterConfig.Filter(ctx)
		if err != nil {
			return nil, err
		}
		filters = append(filters, filter)
	}
	return source.NewFilterSource(src, filters), nil
}

func (f *FilterConfig) Filter(ctx *domain.TemplateContext) (domain.Filter, error) {
	if f.Template != nil {
		return filter.NewTemplateFilter(*f.Template, ctx), nil
	}
	if f.DOM != nil {
		selector, err := f.DOM.Evaluate(ctx.Snapshot())
		if err != nil {
			return nil, err
		}
		return filter.NewDOMFilter(selector), nil
	}
	if f.JSONTransform != nil {
		return filter.NewJSONTransformFilter(f.JSONTransform.Transform, ctx.Snapshot()), nil
	}
	if f.JSONObject != nil {
		return filter.NewJSONObjectFilter(), nil
	}
	if f.JSONArray != nil {
		return filter.NewJSONArrayFilter(ctx.Snapshot(), f.JSONArray.Condition), nil
	}
	if f.JSONArraySort != nil {
		return filter.NewJSONArraySortFilter(f.JSONArraySort.By), nil
	}
	return nil, errors.New("no filters defined")
}
