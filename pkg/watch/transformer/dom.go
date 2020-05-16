package transformer

import (
	"fmt"
	"github.com/uphy/watch-web/pkg/domain/template"
	"github.com/uphy/watch-web/pkg/domain/value"

	"github.com/uphy/watch-web/pkg/domain"
)

type (
	DOMTransformer struct {
		selecter string
	}
)

func NewDOMTransformer(selector string) *DOMTransformer {
	return &DOMTransformer{selector}
}

func (t *DOMTransformer) Transform(ctx *domain.JobContext, v value.Value) (value.Value, error) {
	return template.SelectDOM(v.String(), t.selecter)
}

func (t *DOMTransformer) String() string {
	return fmt.Sprintf("DOM[selector=%v]", t.selecter)
}
