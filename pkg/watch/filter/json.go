package filter

import (
	"fmt"
	"strconv"

	"github.com/uphy/watch-web/pkg/domain"
)

type (
	JSONArrayFilter struct {
		condition *domain.TemplateString
		ctx       *domain.TemplateContext
	}
	JSONObjectFilter struct {
	}
)

func NewJSONArrayFilter(ctx *domain.TemplateContext, condition *domain.TemplateString) *JSONArrayFilter {
	return &JSONArrayFilter{condition, ctx}
}

func (j *JSONArrayFilter) Filter(ctx *domain.JobContext, v domain.Value) (domain.Value, error) {
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
		return domain.NewJSONArrayValue(result), nil
	}
	return domain.ConvertAs(v.String(), domain.ValueTypeJSONArray)
}

func (j *JSONArrayFilter) String() string {
	return fmt.Sprintf("JSONArray[]")
}

func NewJSONObjectFilter() *JSONObjectFilter {
	return &JSONObjectFilter{}
}

func (j *JSONObjectFilter) Filter(ctx *domain.JobContext, v domain.Value) (domain.Value, error) {
	return domain.ConvertAs(v.String(), domain.ValueTypeJSONObject)
}

func (j *JSONObjectFilter) String() string {
	return fmt.Sprintf("JSONObject[]")
}
