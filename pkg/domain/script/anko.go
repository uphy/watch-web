package script

import (
	"fmt"
	"regexp"

	"github.com/mattn/anko/core"
	"github.com/mattn/anko/env"
	"github.com/mattn/anko/parser"
	"github.com/mattn/anko/vm"
	"github.com/uphy/watch-web/pkg/domain"
	"gopkg.in/yaml.v2"
)

type (
	AnkoScriptEngine struct {
	}
	AnkoScript struct {
		e      *env.Env
		script string
	}
)

func NewAnkoScriptEngine() *AnkoScriptEngine {
	return &AnkoScriptEngine{}
}

func (a *AnkoScriptEngine) NewScript(script string) (domain.Script, error) {
	e := env.NewEnv()
	core.Import(e)
	funcs := map[string]interface{}{
		"sprintf": fmt.Sprintf,
		"selectDOM": func(html, selector string) (interface{}, error) {
			v, err := domain.SelectDOM(html, selector)
			if err != nil {
				return nil, err
			}
			return v.Interface(), nil
		},
		"regexReplace": func(s, regex, replacement string) (string, error) {
			r, err := regexp.Compile(regex)
			if err != nil {
				return "", err
			}
			return r.ReplaceAllString(s, replacement), nil
		},
		"printYAML": func(v interface{}) error {
			b, err := yaml.Marshal(v)
			if err != nil {
				return err
			}
			fmt.Println(string(b))
			return nil
		},
	}
	for name, f := range funcs {
		if err := e.Define(name, f); err != nil {
			return nil, fmt.Errorf("Cannot define function: func=%s, err=%w", name, err)
		}
	}
	return &AnkoScript{e, script}, nil
}

func (s *AnkoScript) Evaluate(args map[string]interface{}) (interface{}, error) {
	for k, v := range args {
		if err := s.e.Define(k, v); err != nil {
			return nil, fmt.Errorf("failed to set source to anko script engine:%w", err)
		}
	}

	result, err := vm.Execute(s.e, nil, s.script)
	if err != nil {
		if e, ok := err.(*vm.Error); ok {
			return nil, fmt.Errorf("failed to execute script: line=%d, col=%d, err=%w", e.Pos.Line, e.Pos.Column, e)
		} else if e, ok := err.(*parser.Error); ok {
			return nil, fmt.Errorf("failed to parse script: line=%d, col=%d, err=%w", e.Pos.Line, e.Pos.Column, e)
		}
		return nil, fmt.Errorf("failed to evaluate anko script:%w", err)
	}
	return result, nil
}
