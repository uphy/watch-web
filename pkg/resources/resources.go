package resources

import (
	_ "embed"
)

//go:embed templates/slack.json
var SlackTemplate string

//go:embed templates/slack-array-add.json
var SlackArrayTemplateAdd string

//go:embed templates/slack-array-remove.json
var SlackArrayTemplateRemove string

//go:embed templates/slack-array-change.json
var SlackArrayTemplateChange string
