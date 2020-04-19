package filter

import (
	"fmt"

	"github.com/uphy/watch-web/pkg/domain"
)

type (
	JSONArrayFilter struct {
	}
	JSONObjectFilter struct {
	}
)

func NewJSONArrayFilter() *JSONArrayFilter {
	return &JSONArrayFilter{}
}

func (j *JSONArrayFilter) Filter(ctx *domain.JobContext, v domain.Value) (domain.Value, error) {
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
