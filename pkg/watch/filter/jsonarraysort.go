package filter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/uphy/watch-web/pkg/domain"
)

type (
	JSONArraySortFilter struct {
		by string
	}
)

func NewJSONArraySortFilter(by string) *JSONArraySortFilter {
	return &JSONArraySortFilter{by}
}

func (j *JSONArraySortFilter) Filter(ctx *domain.JobContext, v domain.Value) (domain.Value, error) {
	array := v.JSONArray()
	extract := func(v interface{}) string {
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprint(v)
		}
		var value map[string]interface{}
		if err := json.Unmarshal(b, &value); err != nil {
			return fmt.Sprint(v)
		}
		fieldValue, exist := value[j.by]
		if !exist {
			return fmt.Sprint(v)
		}
		return fmt.Sprint(fieldValue)
	}
	sort.Slice(array, func(i, j int) bool {
		v1 := extract(array[i])
		v2 := extract(array[j])
		return strings.Compare(v1, v2) < 0
	})
	return domain.NewJSONArrayValue(array), nil
}

func (j *JSONArraySortFilter) String() string {
	return fmt.Sprintf("JSONArraySort[by=%s]", j.by)
}
