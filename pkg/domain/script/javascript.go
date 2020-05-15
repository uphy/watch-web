package script

import (
	"fmt"

	"github.com/robertkrimen/otto"
	"github.com/uphy/watch-web/pkg/domain"
)

type (
	JavaScriptEngine struct {
		vm *otto.Otto
	}
	JavaScript struct {
		vm     *otto.Otto
		script *otto.Script
	}
)

func NewJavaScriptEngine() *JavaScriptEngine {
	return &JavaScriptEngine{otto.New()}
}

func (e *JavaScriptEngine) NewScript(script string) (domain.Script, error) {
	s, err := e.vm.Compile("script.js", script)
	if err != nil {
		return nil, fmt.Errorf("failed to compile: %v", err)
	}
	return &JavaScript{e.vm, s}, nil
}

func (s *JavaScript) Evaluate(args map[string]interface{}) (interface{}, error) {
	vm := s.vm.Copy()
	for k, v := range args {
		vm.Set(k, v)
	}
	result, err := vm.Run(s.script)
	if err != nil {
		return nil, err
	}
	exported, err := result.Export()
	if err != nil {
		return nil, err
	}
	return exported, nil
}
