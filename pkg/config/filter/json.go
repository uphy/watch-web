package filter

import (
	"fmt"

	"github.com/uphy/watch-web/pkg/check"
	"github.com/uphy/watch-web/pkg/value"
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

func (j *JSONArrayFilter) Filter(ctx *check.JobContext, v value.Value) (value.Value, error) {
	return value.ConvertAs(v.String(), value.ValueTypeJSONArray)
}

func (j *JSONArrayFilter) String() string {
	return fmt.Sprintf("JSONArray[]")
}

func NewJSONObjectFilter() *JSONObjectFilter {
	return &JSONObjectFilter{}
}

func (j *JSONObjectFilter) Filter(ctx *check.JobContext, v value.Value) (value.Value, error) {
	return value.ConvertAs(v.String(), value.ValueTypeJSONObject)
}

func (j *JSONObjectFilter) String() string {
	return fmt.Sprintf("JSONObject[]")
}
