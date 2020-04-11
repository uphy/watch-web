package resources

import (
	"io/ioutil"
	"log"

	"github.com/markbates/pkger"
)

var SlackTemplate string

func init() {
	f, err := pkger.Open("/resources/slack.json")
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	SlackTemplate = string(b)
}
