package config

import (
	"github.com/uphy/watch-web/pkg/domain"
)

type (
	TransformsConfig []TransformConfig
	TransformConfig  struct {
		Template  *domain.TemplateString `json:"template,omitempty"`
		DOM       *domain.TemplateString `json:"dom,omitempty"`
		JSONArray *struct {
			Condition *domain.TemplateString `json:"condition,omitempty"`
		} `json:"json_array,omitempty"`
		JSONObject *struct{} `json:"json_object,omitempty"`
		Map        *struct {
			Template map[string]domain.TemplateString `json:"template,omitempty"`
		} `json:"map,omitempty"`
		Sort *struct {
			By string `json:"by"`
		} `json:"sort,omitempty"`
		Script *ScriptConfig `json:"script,omitempty"`
		Debug  *bool         `json:"debug"`
	}
)
