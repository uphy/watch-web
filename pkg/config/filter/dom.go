package filter

import (
	"fmt"

	"github.com/uphy/watch-web/pkg/watch"
	"github.com/uphy/watch-web/pkg/config/template"
	"github.com/uphy/watch-web/pkg/value"
)

type (
	DOMFilter struct {
		selecter string
	}
)

func NewDOMFilter(selector string) *DOMFilter {
	return &DOMFilter{selector}
}

func (t *DOMFilter) Filter(ctx *watch.JobContext, v value.Value) (value.Value, error) {
	return template.ParseDOM(v.String(), t.selecter)
}

func (t *DOMFilter) String() string {
	return fmt.Sprintf("DOM[selector=%v]", t.selecter)
}
