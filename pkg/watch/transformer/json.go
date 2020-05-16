package transformer

import (
	"fmt"
	"github.com/uphy/watch-web/pkg/domain/template"
	"github.com/uphy/watch-web/pkg/domain/value"
	"strconv"

	"github.com/uphy/watch-web/pkg/domain"
)

type (
	JSONArrayTransformer struct {
		condition *template.TemplateString
		ctx       *template.TemplateContext
	}
	JSONObjectTransformer struct {
	}
)

func NewJSONArrayTransformer(ctx *template.TemplateContext, condition *template.TemplateString) *JSONArrayTransformer {
	return &JSONArrayTransformer{condition, ctx}
}

func (j *JSONArrayTransformer) Transform(ctx *domain.JobContext, v value.Value) (value.Value, error) {
	if j.condition != nil {
		array, err := value.ConvertAs(v.String(), value.ValueTypeJSONArray)
		if err != nil {
			return nil, err
		}
		result := make(value.JSONArray, 0)
		for _, elm := range array.JSONArray() {
			j.ctx.PushScope()
			j.ctx.Set("source", elm)
			s, err := j.condition.Evaluate(j.ctx)
			j.ctx.PopScope()
			if err != nil {
				return nil, err
			}
			b, _ := strconv.ParseBool(s)
			if b {
				result = append(result, elm)
			}
		}
		return value.NewJSONArray(result), nil
	}
	return value.ConvertAs(v.String(), value.ValueTypeJSONArray)
}

func (j *JSONArrayTransformer) String() string {
	return fmt.Sprintf("JSONArray[]")
}

func NewJSONObjectTransformer() *JSONObjectTransformer {
	return &JSONObjectTransformer{}
}

func (j *JSONObjectTransformer) Transform(ctx *domain.JobContext, v value.Value) (value.Value, error) {
	return value.ConvertAs(v.String(), value.ValueTypeJSONObject)
}

func (j *JSONObjectTransformer) String() string {
	return fmt.Sprintf("JSONObject[]")
}
