package transformer

import (
	"encoding/json"
	"fmt"
	"github.com/uphy/watch-web/pkg/domain/value"
	"sort"
	"strings"

	"github.com/uphy/watch-web/pkg/domain"
)

type (
	SortTransformer struct {
		by string
	}
)

func NewSortTransformer(by string) *SortTransformer {
	return &SortTransformer{by}
}

func (j *SortTransformer) Transform(ctx *domain.JobContext, v value.Value) (value.Value, error) {
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
	return value.NewJSONArray(array), nil
}

func (j *SortTransformer) String() string {
	return fmt.Sprintf("JSONArraySort[by=%s]", j.by)
}
