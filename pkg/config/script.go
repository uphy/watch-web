package config

import (
	"errors"
	"github.com/uphy/watch-web/pkg/domain/template"

	"github.com/uphy/watch-web/pkg/domain"
	"github.com/uphy/watch-web/pkg/domain/script"
)

type (
	ScriptConfig struct {
		JavaScript *template.TemplateString `json:"javascript"`
		Anko       *template.TemplateString `json:"anko"`
	}
)

func (s *ScriptConfig) NewScript(ctx *template.TemplateContext) (domain.Script, error) {
	if s.JavaScript != nil {
		scr, err := s.JavaScript.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		return script.NewJavaScriptEngine().NewScript(scr)
	}
	if s.Anko != nil {
		scr, err := s.Anko.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		return script.NewAnkoScriptEngine().NewScript(scr)
	}
	return nil, errors.New("no script engine defined")
}
