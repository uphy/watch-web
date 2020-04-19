package source

import (
	"fmt"

	"github.com/uphy/watch-web/pkg/domain"
)

type (
	ConstantSource struct {
		value interface{}
	}
)

func NewConstantSource(constant interface{}) *ConstantSource {
	return &ConstantSource{
		value: constant,
	}
}

func (c *ConstantSource) Fetch(ctx *domain.JobContext) (domain.Value, error) {
	return domain.ConvertInterfaceAs(c.value, domain.ValueTypeAutoDetect)
}

func (c *ConstantSource) String() string {
	return fmt.Sprintf("Constant[value=%v]", c.value)
}
