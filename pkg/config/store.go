package config

import (
	"github.com/uphy/watch-web/pkg/domain"
)

type (
	StoreConfig struct {
		Redis *RedisConfig `json:"redis,omitempty"`
	}
	RedisConfig struct {
		Address   *domain.TemplateString `json:"address"`
		Password  *domain.TemplateString `json:"password"`
		RedisToGo *domain.TemplateString `json:"redistogo"`
	}
)
