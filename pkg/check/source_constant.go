package check

import (
	"fmt"

	"github.com/uphy/watch-web/pkg/value"
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

func (c *ConstantSource) Fetch(ctx *JobContext) (value.Value, error) {
	return value.ConvertInterfaceAs(c.value, value.ValueTypeAutoDetect)
}

func (c *ConstantSource) String() string {
	return fmt.Sprintf("Constant[value=%v]", c.value)
}
