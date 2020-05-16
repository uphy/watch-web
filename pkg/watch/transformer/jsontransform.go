package transformer

import (
	"fmt"
	"github.com/uphy/watch-web/pkg/domain/template"
	"github.com/uphy/watch-web/pkg/domain/value"

	"github.com/uphy/watch-web/pkg/domain"
)

type (
	transform      map[string]template.TemplateString
	MapTransformer struct {
		transform transform
		ctx       *template.TemplateContext
	}
)

func NewMapTransformer(transform map[string]template.TemplateString, ctx *template.TemplateContext) *MapTransformer {
	return &MapTransformer{transform, ctx}
}

func (t transform) transform(ctx *template.TemplateContext, v value.Value) (value.Value, error) {
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
	return value.NewJSONObject(transformed), nil
}

func (t *MapTransformer) Transform(ctx *domain.JobContext, v value.Value) (value.Value, error) {
	switch v.Type() {
	case value.ValueTypeString, value.ValueTypeJSONObject:
		return t.transform.transform(t.ctx, v)
	case value.ValueTypeJSONArray:
		var result = make([]interface{}, 0)
		for _, elm := range v.JSONArray() {
			value, err := value.ConvertInterfaceAs(elm, value.ValueTypeJSONObject)
			if err != nil {
				return nil, err
			}
			transformed, err := t.transform.transform(t.ctx, value)
			if err != nil {
				return nil, err
			}
			result = append(result, transformed.Interface())
		}
		return value.NewJSONArray(result), nil
	default:
		return nil, fmt.Errorf("unsupported value type: %s", v.Type())
	}
}
