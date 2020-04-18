package resources

import (
	"io/ioutil"
	"log"

	"github.com/markbates/pkger"
)

var SlackTemplate string
var HttpStatic pkger.Dir

func init() {
	f, err := pkger.Open("/templates/slack.json")
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	HttpStatic = pkger.Dir("/frontend/dist")
	SlackTemplate = string(b)
}
