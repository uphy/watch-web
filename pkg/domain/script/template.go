package script

import (
	"github.com/uphy/watch-web/pkg/domain"
	"github.com/uphy/watch-web/pkg/domain/template"
)

type (
	TemplateScriptEngine struct {
		ctx *template.TemplateContext
	}
	TemplateScript struct {
		ctx      *template.TemplateContext
		template *template.Template
	}
)

func NewTemplateScriptEngine(ctx *template.TemplateContext) *TemplateScriptEngine {
	return &TemplateScriptEngine{ctx}
}

func (t *TemplateScriptEngine) NewScript(script string) (domain.Script, error) {
	tmpl, err := template.Parse(script)
	if err != nil {
		return nil, err
	}
	return &TemplateScript{t.ctx, tmpl}, nil
}

func (s *TemplateScript) Evaluate(args map[string]interface{}) (interface{}, error) {
	s.ctx.PushScope()
	defer s.ctx.PopScope()
	for k, v := range args {
		s.ctx.Set(k, v)
	}
	return s.template.Evaluate(s.ctx)
}
