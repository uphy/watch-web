package config

import (
	"github.com/uphy/watch-web/pkg/domain"
)

type (
	JobConfig struct {
		ID        domain.TemplateString `json:"id"`
		Label     domain.TemplateString `json:"label"`
		Link      domain.TemplateString `json:"link"`
		Source    *SourceConfig         `json:"source,omitempty"`
		Schedule  domain.TemplateString `json:"schedule,omitempty"`
		WithItems []interface{}         `json:"with_items,omitempty"`
	}
)
