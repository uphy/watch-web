package action

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
	"github.com/uphy/watch-web/pkg/domain"
)

type (
	SlackBotActionRepository interface {
		PutSlackThreadTS(jobID string, itemID string, threadTS string) error
		GetSlackThreadTS(jobID, itemID string) (*string, error)
	}
	SlackBotAction struct {
		token   string
		channel string
		debug   bool
		repo    SlackBotActionRepository
	}
	slackBotResponse struct {
		OK      bool   `json:"ok"`
		Channel string `json:"channel"`
		TS      string `json:"ts"`
	}
)

func NewSlackBotAction(token string, channel string, debug bool, repo SlackBotActionRepository) *SlackBotAction {
	return &SlackBotAction{token, channel, debug, repo}
}

func (s *SlackBotAction) Run(ctx *domain.JobContext, res *domain.Result) error {
	updates := res.Diff()
	if !updates.Changes() {
		return nil
	}
	for _, update := range updates {
		payloadValue, err := slackPayload(ctx, res, update)
		if err != nil {
			return err
		}
		if payloadValue == nil {
			return nil
		}

		// Set additional request parameters
		payloadValue["channel"] = s.channel
		payloadValue["unfurl_links"] = false
		payloadValue["unfurl_media"] = false
		isReply := false
		if s.repo != nil {
			threadTS, err := s.repo.GetSlackThreadTS(res.JobID, update.ItemID())
			if err != nil {
				ctx.Log.WithFields(logrus.Fields{
					"jobID":  res.JobID,
					"itemID": update.ItemID(),
				}).Warn("failed to get slack thread_ts from repository.  Post to channel.")
			}
			if threadTS != nil {
				payloadValue["thread_ts"] = *threadTS
				payloadValue["reply_broadcast"] = true
				isReply = true
			}
		}
		if s.debug || len(s.token) == 0 {
			ctx.Log.Info("Slack action is debug mode.  No notification.")
			payloadBytes, _ := yaml.Marshal(payloadValue)
			fmt.Println("[Payload]")
			fmt.Println(string(payloadBytes))
		} else {
			payloadBytes, err := json.Marshal(payloadValue)
			if err != nil {
				return err
			}
			post, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewReader(payloadBytes))
			if err != nil {
				return err
			}
			post.Header.Add("Content-Type", "application/json; charset=UTF-8")
			post.Header.Add("Authorization", "Bearer "+s.token)
			resp, err := http.DefaultClient.Do(post)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			b, _ := ioutil.ReadAll(resp.Body)
			if resp.StatusCode != 200 {
				return fmt.Errorf("invalid status code: status=%d, body=%s", resp.StatusCode, string(b))
			}
			var response slackBotResponse
			if err := json.Unmarshal(b, &response); err != nil {
				return fmt.Errorf("failed to unmarshal slack api response: err=%v, body=%s", err, string(b))
			}
			threadTS := response.TS
			if len(threadTS) > 0 && s.repo != nil && !isReply { // if this message is posted to thread, then no need to save timestamp
				if err := s.repo.PutSlackThreadTS(res.JobID, update.ItemID(), threadTS); err != nil {
					ctx.Log.WithFields(logrus.Fields{
						"jobID":  res.JobID,
						"itemID": update.ItemID(),
					}).Warn("failed to put slack thread_ts into repository.")
				}
			}
		}
	}
	return nil
}
