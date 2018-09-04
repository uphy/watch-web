package check

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type (
	slackPayload struct {
		Text string `json:"text"`
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
	payload := &slackPayload{
		Text: fmt.Sprintf(`Updated %s

%s
%s
`, res.Name, changes.String(), res.Label),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := http.Post(s.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("invalid status code: " + strconv.Itoa(resp.StatusCode))
	}
	return nil
}
