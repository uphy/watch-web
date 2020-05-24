package resources

import (
	"io/ioutil"
	"log"

	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
)

var SlackTemplate string
var SlackArrayTemplateAdd string
var SlackArrayTemplateRemove string
var SlackArrayTemplateChange string

func init() {
	SlackTemplate = load(pkger.Open("/templates/slack.json"))
	SlackArrayTemplateAdd = load(pkger.Open("/templates/slack-array-add.json"))
	SlackArrayTemplateRemove = load(pkger.Open("/templates/slack-array-remove.json"))
	SlackArrayTemplateChange = load(pkger.Open("/templates/slack-array-change.json"))
}

func load(f pkging.File, err error) string {
	if err != nil {
		log.Fatalf("failed to open embedded template file: err=%v", err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}
