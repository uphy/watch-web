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
		JSONObject    *struct{} `json:"json_object,omitempty"`
		JSONTransform *struct {
			Transform map[string]domain.TemplateString `json:"transform,omitempty"`
		} `json:"json_transform,omitempty"`
		JSONArraySort *struct {
			By string `json:"by"`
		} `json:"json_array_sort"`
		Script *struct {
			Script   *domain.TemplateString `json:"script"`
			Language *string                `json:"lang"`
		} `json:"script,omitempty"`
		Debug *bool `json:"debug"`
	}
)
