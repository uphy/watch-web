package transformer

import (
	"fmt"
	"github.com/uphy/watch-web/pkg/domain/value"

	"github.com/ghodss/yaml"
	"github.com/uphy/watch-web/pkg/domain"
)

type (
	DebugTransformer struct {
		Debug bool
	}
)

func NewDebugTransformer(debug bool) *DebugTransformer {
	return &DebugTransformer{debug}
}

func (t DebugTransformer) Transform(ctx *domain.JobContext, v value.Value) (value.Value, error) {
	b, err := yaml.Marshal(v.Interface())
	if err != nil {
		return nil, err
	}
	fmt.Printf("[Type]\n%s\n", v.Type())
	fmt.Println("[Value]")
	fmt.Println(string(b))
	return v, nil
}
