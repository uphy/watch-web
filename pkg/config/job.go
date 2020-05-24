package config

import (
	"github.com/uphy/watch-web/pkg/domain/template"
)

type (
	JobConfig struct {
		ID        template.TemplateString `json:"id"`
		Label     template.TemplateString `json:"label"`
		Link      template.TemplateString `json:"link"`
		Source    *SourceConfig           `json:"source,omitempty"`
		Schedule  template.TemplateString `json:"schedule,omitempty"`
		WithItems []interface{}           `json:"with_items,omitempty"`
		Actions   []ActionConfig          `json:"actions",omitempty`
	}
)
