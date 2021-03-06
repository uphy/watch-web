package action

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/uphy/watch-web/pkg/domain/template"
	"github.com/uphy/watch-web/pkg/domain/value"

	"github.com/uphy/watch-web/pkg/domain"
	"github.com/uphy/watch-web/pkg/resources"
)

type (
	SlackWebhookAction struct {
		URL   string
		Debug bool
	}
)

func NewSlackWebhookAction(webhookURL string, debug bool) *SlackWebhookAction {
	return &SlackWebhookAction{webhookURL, debug}
}

func (s *SlackWebhookAction) Run(ctx *domain.JobContext, res *domain.Result) error {
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

		if s.Debug || len(s.URL) == 0 {
			ctx.Log.Info("Slack action is debug mode.  No notification.")
			payloadBytes, _ := yaml.Marshal(payloadValue)
			fmt.Println("[Payload]")
			fmt.Println(string(payloadBytes))
		} else {
			payloadBytes, err := json.Marshal(payloadValue)
			if err != nil {
				return err
			}
			resp, err := http.Post(s.URL, "application/json", bytes.NewReader(payloadBytes))
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			b, _ := ioutil.ReadAll(resp.Body)
			if resp.StatusCode != 200 {
				return fmt.Errorf("invalid status code: status=%d, body=%s", resp.StatusCode, string(b))
			}
		}
	}
	return nil
}

func slackPayload(ctx *domain.JobContext, res *domain.Result, update value.Update) (map[string]interface{}, error) {
	var tmpl *template.Template
	args := map[string]interface{}{
		"res": res,
	}
	var err error
	switch item := update.Item().(type) {
	// add or remove
	case value.Item:
		if update.Type == value.UpdateTypeAdd {
			tmpl, err = template.Parse(resources.SlackArrayTemplateAdd)
		} else {
			tmpl, err = template.Parse(resources.SlackArrayTemplateRemove)
		}
		fields := make([]map[string]string, 0)
		for k, v := range item {
			if k == "link" || k == "summary" || k == "thumbnail" || k == "id" {
				continue
			}
			fields = append(fields, map[string]string{
				"key":   k,
				"value": v,
			})
		}
		sort.Slice(fields, func(i, j int) bool {
			k1 := fields[i]["key"]
			k2 := fields[j]["key"]
			return strings.Compare(k1, k2) < 0
		})
		args["fields"] = fields
		args["item"] = update.Item().(value.Item)
	// change
	case *value.ItemChange:
		tmpl, err = template.Parse(resources.SlackArrayTemplateChange)
		fields := make([]map[string]string, 0)
		for k, v := range item.AddedKeys {
			fields = append(fields, map[string]string{
				"type":  "add",
				"key":   k,
				"value": v,
			})
		}
		for k, v := range item.ChangedKeys {
			fields = append(fields, map[string]string{
				"type": "change",
				"key":  k,
				"old":  v.Old,
				"new":  v.New,
			})
		}
		for k, v := range item.RemovedKeys {
			fields = append(fields, map[string]string{
				"type":  "remove",
				"key":   k,
				"value": v,
			})
		}
		sort.Slice(fields, func(i, j int) bool {
			k1 := fields[i]["key"]
			k2 := fields[j]["key"]
			return strings.Compare(k1, k2) < 0
		})
		args["fields"] = fields
		args["item"] = update.Change.Item
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %v", err)
	}

	templateContext := template.NewRootTemplateContext()
	for k, v := range args {
		templateContext.Set(k, v)
	}
	payload, err := tmpl.Evaluate(templateContext)
	if err != nil {
		return nil, err
	}

	var payloadValue map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &payloadValue); err != nil {
		ctx.Log.WithField("err", err).Error("Failed to parse payload as json")
		fmt.Println(payload)
		return nil, err
	}
	return payloadValue, nil
}
