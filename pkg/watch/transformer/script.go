package transformer

import (
	"fmt"
	"github.com/uphy/watch-web/pkg/domain/value"

	"github.com/uphy/watch-web/pkg/domain"
)

type (
	ScriptTransformer struct {
		script domain.Script
	}
)

func NewScriptTransformer(script domain.Script) (*ScriptTransformer, error) {
	return &ScriptTransformer{script}, nil
}

func (t ScriptTransformer) Transform(ctx *domain.JobContext, v value.Value) (value.Value, error) {
	result, err := t.script.Evaluate(map[string]interface{}{
		"source": v,
	})

	if err != nil {
		return nil, err
	}

	switch res := result.(type) {
	case value.Value:
		return res, nil
	case map[string]interface{}:
		return value.NewJSONObject(res), nil
	case map[string]string:
		m := make(map[string]interface{})
		for key, value := range res {
			m[key] = value
		}
		return value.NewJSONObject(m), nil
	case string:
		return value.NewStringValue(res), nil
	}
	return nil, fmt.Errorf("Unsupported script return value: %v", result)
}
