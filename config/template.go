package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type (
	TemplateString  string
	TemplateContext struct {
		scope *templateScope
	}
	templateScope struct {
		parent    *templateScope
		variables map[string]interface{}
	}
)

var funcs = map[string]interface{}{
	"env": func(name string) string {
		return os.Getenv(name)
	},
	"default": func(defaultValue string, value string) string {
		if value == "" {
			return defaultValue
		}
		return value
	},
	"sliceOf": func(values ...string) []string {
		return values
	},
	"json": func(jsonString string) (interface{}, error) {
		if strings.Trim(jsonString, " ") == "" {
			return make(map[string]interface{}), nil
		}
		var v interface{}
		if err := json.Unmarshal([]byte(jsonString), &v); err != nil {
			return nil, err
		}
		return v, nil
	},
	"dom": func(selector, html string) (interface{}, error) {
		// this returns parsed DOM as map(json object value)
		v, err := parseDOM(html, selector)
		if err != nil {
			return nil, err
		}
		return v.Interface(), nil
	},
}

func (t TemplateString) Evaluate(ctx *TemplateContext) (string, error) {
	var buf = new(bytes.Buffer)
	if err := template.Must(template.New("template-string").Funcs(funcs).Parse(string(t))).Execute(buf, ctx.Get()); err != nil {
		return "", fmt.Errorf("cannot evaluate %s: %w", t, err)
	}
	return buf.String(), nil
}

func NewRootContext() *TemplateContext {
	return &TemplateContext{&templateScope{nil, make(map[string]interface{})}}
}

func (c *TemplateContext) PushScope() {
	c.scope = c.scope.child()
}

func (c *TemplateContext) PopScope() {
	c.scope = c.scope.parent
}

func (c *TemplateContext) Get() map[string]interface{} {
	return c.scope.get()
}

func (c *TemplateContext) Set(key string, value interface{}) {
	c.scope.set(key, value)
}

func (c *templateScope) child() *templateScope {
	return &templateScope{c, make(map[string]interface{})}
}

func (c *templateScope) set(key string, value interface{}) {
	c.variables[key] = value
}

func (c *templateScope) get() map[string]interface{} {
	if c.parent != nil {
		vars := make(map[string]interface{})
		parentVars := c.parent.get()
		for k, v := range parentVars {
			vars[k] = v
		}
		for k, v := range c.variables {
			vars[k] = v
		}
		return vars
	}
	return c.variables
}
