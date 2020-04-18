package filter

import (
	"fmt"
	"strings"

	"github.com/uphy/watch-web/pkg/check"
	"github.com/uphy/watch-web/pkg/config/template"
	"github.com/uphy/watch-web/pkg/value"
)

const (
	JSONTransformSourceTypeAuto   = "auto"
	JSONTransformSourceTypeArray  = "array"
	JSONTransformSourceTypeObject = "object"
	JSONTransformSourceTypeRaw    = "raw"
)

type (
	JSONTransformSourceType string
	JSONTransformFilter     struct {
		sourceType JSONTransformSourceType
		transform  map[string]template.TemplateString
		ctx        *template.TemplateContext
	}
)

func NewJSONTransformFilter(sourceType JSONTransformSourceType, transform map[string]template.TemplateString, ctx *template.TemplateContext) *JSONTransformFilter {
	return &JSONTransformFilter{sourceType, transform, ctx}
}

func (t *JSONTransformFilter) Filter(ctx *check.JobContext, v value.Value) (value.Value, error) {
	// auto detect source type
	var sourceType = t.sourceType
	if sourceType == JSONTransformSourceTypeAuto {
		s := v.String()
		if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
			sourceType = JSONTransformSourceTypeArray
		} else if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
			sourceType = JSONTransformSourceTypeObject
		} else {
			sourceType = JSONTransformSourceTypeRaw
		}
	}
	// parse source based on type
	var source []interface{}
	hasMultiElements := false
	switch sourceType {
	case JSONTransformSourceTypeArray:
		source = v.JSONArray()
		hasMultiElements = true
	case JSONTransformSourceTypeObject:
		source = []interface{}{v.JSONObject()}
	case JSONTransformSourceTypeRaw:
		source = []interface{}{v.Interface()}
	default:
		return nil, fmt.Errorf("unsupported transform source type: %v", sourceType)
	}
	// transform
	var transformed []interface{}
	for _, src := range source {
		t.ctx.PushScope()
		t.ctx.Set("source", src)
		elm := make(map[string]string)
		for k, tmpl := range t.transform {
			evaluated, err := tmpl.Evaluate(t.ctx)
			if err != nil {
				t.ctx.PopScope()
				return nil, err
			}
			elm[k] = evaluated
		}
		transformed = append(transformed, elm)
		t.ctx.PopScope()
	}
	// return
	if hasMultiElements {
		return value.NewJSONArrayValue(transformed), nil
	}
	return value.NewJSONObjectValue(transformed[0].(map[string]interface{})), nil
}
