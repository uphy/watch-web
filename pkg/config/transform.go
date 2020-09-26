package config

import (
	"github.com/uphy/watch-web/pkg/domain/template"
)

type (
	TransformsConfig []TransformConfig
	TransformConfig  struct {
		Template  *template.TemplateString `json:"template,omitempty"`
		DOM       *template.TemplateString `json:"dom,omitempty"`
		JSONArray *struct {
			Condition *template.TemplateString `json:"condition,omitempty"`
		} `json:"json_array,omitempty"`
		JSONObject *struct{} `json:"json_object,omitempty"`
		Map        *struct {
			Template map[string]template.TemplateString `json:"template,omitempty"`
		} `json:"map,omitempty"`
		Sort *struct {
			By string `json:"by"`
		} `json:"sort,omitempty"`
		Script *ScriptConfig `json:"script,omitempty"`
		Filter *ScriptConfig `json:"filter,omitempty"`
		Debug  *bool         `json:"debug"`
		Retry  interface{}   `json:"retry"`
	}
)
