package filter

import (
	"fmt"

	"github.com/uphy/watch-web/pkg/domain"
)

type (
	TemplateFilter struct {
		template domain.TemplateString
		ctx      *domain.TemplateContext
	}
)

func NewTemplateFilter(template domain.TemplateString, ctx *domain.TemplateContext) *TemplateFilter {
	return &TemplateFilter{template, ctx}
}

func (t *TemplateFilter) Filter(ctx *domain.JobContext, v domain.Value) (domain.Value, error) {
	t.ctx.PushScope()
	defer t.ctx.PopScope()
	t.ctx.Set("source", v.Interface())
	evaluated, err := t.template.Evaluate(t.ctx)
	if err != nil {
		return nil, err
	}
	return domain.Auto(evaluated), nil
}

func (t *TemplateFilter) String() string {
	return fmt.Sprintf("Template[template=%v]", t.template)
}
