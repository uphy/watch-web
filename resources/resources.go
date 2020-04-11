package resources

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/markbates/pkger"
)

var SlackTemplate string
var HttpStatic pkger.Dir

func init() {
	f, err := pkger.Open("/resources/slack.json")
	if err != nil {
		pkger.Walk("/", func(path string, info os.FileInfo, err error) error {
			fmt.Println(path)
			return nil
		})
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	HttpStatic = pkger.Dir("/frontend/dist")
	SlackTemplate = string(b)
}
