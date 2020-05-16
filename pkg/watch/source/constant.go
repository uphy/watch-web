package source

import (
	"fmt"
	value2 "github.com/uphy/watch-web/pkg/domain/value"

	"github.com/uphy/watch-web/pkg/domain"
)

type (
	ConstantSource struct {
		value value2.Value
	}
)

func NewConstantSource(constant value2.Value) *ConstantSource {
	return &ConstantSource{
		value: constant,
	}
}

func (c *ConstantSource) Fetch(ctx *domain.JobContext) (value2.Value, error) {
	return c.value, nil
}

func (c *ConstantSource) String() string {
	return fmt.Sprintf("Constant[value=%v]", c.value)
}
