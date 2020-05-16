package config

import (
	"github.com/uphy/watch-web/pkg/domain/template"
)

type (
	StoreConfig struct {
		Redis *RedisConfig `json:"redis,omitempty"`
	}
	RedisConfig struct {
		Address   *template.TemplateString `json:"address"`
		Password  *template.TemplateString `json:"password"`
		RedisToGo *template.TemplateString `json:"redistogo"`
	}
)
