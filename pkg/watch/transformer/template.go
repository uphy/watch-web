package transformer

import (
	"fmt"

	"github.com/uphy/watch-web/pkg/domain"
)

type (
	TemplateTransformer struct {
		template domain.TemplateString
		ctx      *domain.TemplateContext
	}
)

func NewTemplateTransformer(template domain.TemplateString, ctx *domain.TemplateContext) *TemplateTransformer {
	return &TemplateTransformer{template, ctx}
}

func (t *TemplateTransformer) Transform(ctx *domain.JobContext, v domain.Value) (domain.Value, error) {
	t.ctx.PushScope()
	defer t.ctx.PopScope()
	t.ctx.Set("source", v.Interface())
	evaluated, err := t.template.Evaluate(t.ctx)
	if err != nil {
		return nil, err
	}
	return domain.Auto(evaluated), nil
}

func (t *TemplateTransformer) String() string {
	return fmt.Sprintf("Template[template=%v]", t.template)
}
