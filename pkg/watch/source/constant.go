package source

import (
	"fmt"

	"github.com/uphy/watch-web/pkg/domain"
)

type (
	ConstantSource struct {
		value domain.Value
	}
)

func NewConstantSource(constant domain.Value) *ConstantSource {
	return &ConstantSource{
		value: constant,
	}
}

func (c *ConstantSource) Fetch(ctx *domain.JobContext) (domain.Value, error) {
	return c.value, nil
}

func (c *ConstantSource) String() string {
	return fmt.Sprintf("Constant[value=%v]", c.value)
}
