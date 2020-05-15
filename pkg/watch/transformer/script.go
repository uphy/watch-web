package transformer

import (
	"fmt"

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

func (t ScriptTransformer) Transform(ctx *domain.JobContext, v domain.Value) (domain.Value, error) {
	result, err := t.script.Evaluate(map[string]interface{}{
		"source": v,
	})

	if err != nil {
		return nil, err
	}

	switch res := result.(type) {
	case domain.Value:
		return res, nil
	case map[string]interface{}:
		return domain.NewJSONObject(res), nil
	case map[string]string:
		m := make(map[string]interface{})
		for key, value := range res {
			m[key] = value
		}
		return domain.NewJSONObject(m), nil
	case string:
		return domain.NewStringValue(res), nil
	}
	return nil, fmt.Errorf("Unsupported script return value: %v", result)
}
