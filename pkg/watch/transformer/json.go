package transformer

import (
	"fmt"
	"strconv"

	"github.com/uphy/watch-web/pkg/domain"
)

type (
	JSONArrayTransformer struct {
		condition *domain.TemplateString
		ctx       *domain.TemplateContext
	}
	JSONObjectTransformer struct {
	}
)

func NewJSONArrayTransformer(ctx *domain.TemplateContext, condition *domain.TemplateString) *JSONArrayTransformer {
	return &JSONArrayTransformer{condition, ctx}
}

func (j *JSONArrayTransformer) Transform(ctx *domain.JobContext, v domain.Value) (domain.Value, error) {
	if j.condition != nil {
		array, err := domain.ConvertAs(v.String(), domain.ValueTypeJSONArray)
		if err != nil {
			return nil, err
		}
		result := make(domain.JSONArray, 0)
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
		return domain.NewJSONArray(result), nil
	}
	return domain.ConvertAs(v.String(), domain.ValueTypeJSONArray)
}

func (j *JSONArrayTransformer) String() string {
	return fmt.Sprintf("JSONArray[]")
}

func NewJSONObjectTransformer() *JSONObjectTransformer {
	return &JSONObjectTransformer{}
}

func (j *JSONObjectTransformer) Transform(ctx *domain.JobContext, v domain.Value) (domain.Value, error) {
	return domain.ConvertAs(v.String(), domain.ValueTypeJSONObject)
}

func (j *JSONObjectTransformer) String() string {
	return fmt.Sprintf("JSONObject[]")
}
