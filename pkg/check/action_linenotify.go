package check

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/uphy/watch-web/pkg/result"
)

type (
	lineNotifyPayload struct {
		Text string `json:"text"`
	}
	LINENotifyAction struct {
		AccessToken string
	}
)

func NewLINENotifyAction(accessToken string) *LINENotifyAction {
	return &LINENotifyAction{accessToken}
}

func (a *LINENotifyAction) Run(res *result.Result) error {
	changes := res.Diff()
	if !changes.Changed() {
		return nil
	}
	message := fmt.Sprintf(`%s に更新があったよ！

%s
%s
`, res.JobID, changes.String(), res.Label)
	form := url.Values{}
	form.Add("message", message)
	req, err := http.NewRequest("POST", "https://notify-api.line.me/api/notify", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+a.AccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("failed to post message: " + resp.Status)
	}
	return nil
}
