package resources

import (
	"io/ioutil"
	"log"

	"github.com/markbates/pkger"
)

var SlackTemplate string
var SlackArrayTemplate string
var HttpStatic pkger.Dir

func init() {
	{
		f, err := pkger.Open("/templates/slack.json")
		if err != nil {
			log.Fatal(err)
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}
		SlackTemplate = string(b)
	}
	{
		f, err := pkger.Open("/templates/slack-array.json")
		if err != nil {
			log.Fatal(err)
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}
		SlackArrayTemplate = string(b)
	}

	HttpStatic = pkger.Dir("/frontend/dist")

}
