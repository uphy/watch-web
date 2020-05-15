package transformer

import (
	"fmt"

	"github.com/uphy/watch-web/pkg/domain"
)

type (
	transform                map[string]domain.TemplateString
	JSONTransformTransformer struct {
		transform transform
		ctx       *domain.TemplateContext
	}
)

func NewJSONTransformTransformer(transform map[string]domain.TemplateString, ctx *domain.TemplateContext) *JSONTransformTransformer {
	return &JSONTransformTransformer{transform, ctx}
}

func (t transform) transform(ctx *domain.TemplateContext, v domain.Value) (domain.Value, error) {
	ctx.PushScope()
	ctx.Set("source", v.Interface())
	defer ctx.PopScope()

	var transformed = make(map[string]interface{})
	for k, tmpl := range t {
		evaluated, err := tmpl.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		transformed[k] = evaluated
	}
	return domain.NewJSONObject(transformed), nil
}

func (t *JSONTransformTransformer) Transform(ctx *domain.JobContext, v domain.Value) (domain.Value, error) {
	switch v.Type() {
	case domain.ValueTypeString, domain.ValueTypeJSONObject:
		return t.transform.transform(t.ctx, v)
	case domain.ValueTypeJSONArray:
		var result = make([]interface{}, 0)
		for _, elm := range v.JSONArray() {
			value, err := domain.ConvertInterfaceAs(elm, domain.ValueTypeJSONObject)
			if err != nil {
				return nil, err
			}
			transformed, err := t.transform.transform(t.ctx, value)
			if err != nil {
				return nil, err
			}
			result = append(result, transformed.Interface())
		}
		return domain.NewJSONArray(result), nil
	default:
		return nil, fmt.Errorf("unsupported value type: %s", v.Type())
	}
}
