package transformer

import (
	"fmt"
	template2 "github.com/uphy/watch-web/pkg/domain/template"
	"github.com/uphy/watch-web/pkg/domain/value"

	"github.com/uphy/watch-web/pkg/domain"
)

type (
	TemplateTransformer struct {
		template template2.TemplateString
		ctx      *template2.TemplateContext
	}
)

func NewTemplateTransformer(template template2.TemplateString, ctx *template2.TemplateContext) *TemplateTransformer {
	return &TemplateTransformer{template, ctx}
}

func (t *TemplateTransformer) Transform(ctx *domain.JobContext, v value.Value) (value.Value, error) {
	t.ctx.PushScope()
	defer t.ctx.PopScope()
	t.ctx.Set("source", v.Interface())
	evaluated, err := t.template.Evaluate(t.ctx)
	if err != nil {
		return nil, err
	}
	return value.Auto(evaluated), nil
}

func (t *TemplateTransformer) String() string {
	return fmt.Sprintf("Template[template=%v]", t.template)
}
