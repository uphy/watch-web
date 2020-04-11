package check

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"
	"text/template"

	"github.com/uphy/watch-web/server/resources"
)

type (
	slackPayload struct {
		Text        string        `json:"text"`
		Blocks      []interface{} `json:"blocks"`
		Attachments []struct {
			Color  string        `json:"color"`
			Blocks []interface{} `json:""`
		}
	}
	SlackAction struct {
		URL string `json:"url"`
	}
)

func NewSlackAction(webhookURL string) *SlackAction {
	return &SlackAction{webhookURL}
}

func (s *SlackAction) Run(res *Result) error {
	changes := res.Diff()
	if !changes.Changed() {
		return nil
	}
	tmpl := template.Must(template.New("slack-template").Parse(resources.SlackTemplate))
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, map[string]interface{}{
		"res":     res,
		"changes": changes,
	}); err != nil {
		return err
	}
	resp, err := http.Post(s.URL, "application/json", bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("invalid status code: " + strconv.Itoa(resp.StatusCode))
	}
	return nil
}
