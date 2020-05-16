package config

import (
	"errors"

	"github.com/uphy/watch-web/pkg/domain"
	"github.com/uphy/watch-web/pkg/domain/script"
)

type (
	ScriptConfig struct {
		JavaScript *domain.TemplateString `json:"javascript"`
		Anko       *domain.TemplateString `json:"anko"`
	}
)

func (s *ScriptConfig) NewScript(ctx *domain.TemplateContext) (domain.Script, error) {
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
