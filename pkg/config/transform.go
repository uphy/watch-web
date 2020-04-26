package config

import (
	"errors"

	"github.com/uphy/watch-web/pkg/domain"
	"github.com/uphy/watch-web/pkg/watch/source"
	"github.com/uphy/watch-web/pkg/watch/transformer"
)

type (
	TransformsConfig []TransformConfig
	TransformConfig  struct {
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

func (t TransformsConfig) Transforms(ctx *domain.TemplateContext, src domain.Source) (domain.Source, error) {
	if len(t) == 0 {
		return src, nil
	}
	transformers := make([]domain.Transformer, 0)
	for _, transformConfig := range t {
		transformer, err := transformConfig.Transform(ctx)
		if err != nil {
			return nil, err
		}
		transformers = append(transformers, transformer)
	}
	return source.NewTransformerSource(src, transformers), nil
}

func (t *TransformConfig) Transform(ctx *domain.TemplateContext) (domain.Transformer, error) {
	if t.Template != nil {
		return transformer.NewTemplateTransformer(*t.Template, ctx), nil
	}
	if t.DOM != nil {
		selector, err := t.DOM.Evaluate(ctx.Snapshot())
		if err != nil {
			return nil, err
		}
		return transformer.NewDOMTransformer(selector), nil
	}
	if t.JSONTransform != nil {
		return transformer.NewJSONTransformTransformer(t.JSONTransform.Transform, ctx.Snapshot()), nil
	}
	if t.JSONObject != nil {
		return transformer.NewJSONObjectTransformer(), nil
	}
	if t.JSONArray != nil {
		return transformer.NewJSONArrayTransformer(ctx.Snapshot(), t.JSONArray.Condition), nil
	}
	if t.JSONArraySort != nil {
		return transformer.NewJSONArraySortTransformer(t.JSONArraySort.By), nil
	}
	return nil, errors.New("no transforms defined")
}
