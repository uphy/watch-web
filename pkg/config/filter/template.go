package filter

import (
	"fmt"

	"github.com/uphy/watch-web/pkg/watch"
	"github.com/uphy/watch-web/pkg/config/template"
	"github.com/uphy/watch-web/pkg/value"
)

type (
	TemplateFilter struct {
		template template.TemplateString
		ctx      *template.TemplateContext
	}
)

func NewTemplateFilter(template template.TemplateString, ctx *template.TemplateContext) *TemplateFilter {
	return &TemplateFilter{template, ctx}
}

func (t *TemplateFilter) Filter(ctx *watch.JobContext, v value.Value) (value.Value, error) {
	t.ctx.PushScope()
	defer t.ctx.PopScope()
	t.ctx.Set("source", v.Interface())
	evaluated, err := t.template.Evaluate(t.ctx)
	if err != nil {
		return nil, err
	}
	return value.Auto(evaluated), nil
}

func (t *TemplateFilter) String() string {
	return fmt.Sprintf("Template[template=%v]", t.template)
}
