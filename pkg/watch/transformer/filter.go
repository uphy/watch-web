package transformer

import (
	"fmt"
	"github.com/uphy/watch-web/pkg/domain/value"

	"github.com/uphy/watch-web/pkg/domain"
)

type FilterTransformer struct {
	script domain.Script
}

func NewFilterTransformer(script domain.Script) *FilterTransformer {
	return &FilterTransformer{script}
}

func (f *FilterTransformer) Transform(ctx *domain.JobContext, v value.Value) (value.Value, error) {
	a := v.JSONArray()
	filtered := make(value.JSONArray, 0)
	for _, elm := range a {
		result, err := f.script.Evaluate(map[string]interface{}{
			"source": elm,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to filter array element due to script evaluation error: %w", err)
		}
		switch r := result.(type) {
		case bool:
			if !r {
				continue
			}
		case int:
			if r == 0 {
				continue
			}
		case float64:
			if r == 0 {
				continue
			}
		}
		filtered = append(filtered, elm)
	}
	return value.NewJSONArray(filtered), nil
}
