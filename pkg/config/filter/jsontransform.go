package filter

import (
	"fmt"

	"github.com/uphy/watch-web/pkg/watch"
	"github.com/uphy/watch-web/pkg/config/template"
	"github.com/uphy/watch-web/pkg/value"
)

type (
	transform           map[string]template.TemplateString
	JSONTransformFilter struct {
		transform transform
		ctx       *template.TemplateContext
	}
)

func NewJSONTransformFilter(transform map[string]template.TemplateString, ctx *template.TemplateContext) *JSONTransformFilter {
	return &JSONTransformFilter{transform, ctx}
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
	return value.NewJSONObjectValue(transformed), nil
}

func (t *JSONTransformFilter) Filter(ctx *watch.JobContext, v value.Value) (value.Value, error) {
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
		return value.NewJSONArrayValue(result), nil
	default:
		return nil, fmt.Errorf("unsupported value type: %s", v.Type())
	}
}
