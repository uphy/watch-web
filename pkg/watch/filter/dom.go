package filter

import (
	"fmt"

	"github.com/uphy/watch-web/pkg/domain"
)

type (
	DOMFilter struct {
		selecter string
	}
)

func NewDOMFilter(selector string) *DOMFilter {
	return &DOMFilter{selector}
}

func (t *DOMFilter) Filter(ctx *domain.JobContext, v domain.Value) (domain.Value, error) {
	return domain.ParseDOM(v.String(), t.selecter)
}

func (t *DOMFilter) String() string {
	return fmt.Sprintf("DOM[selector=%v]", t.selecter)
}
