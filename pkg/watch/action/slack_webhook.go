package action

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
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
	tmpl := template.New("slack-template-" + strings.ToLower(string(update.Type))).Funcs(template.FuncMap{
		"escape": func(s string) string {
			s = strings.ReplaceAll(s, "\n", "\\n")
			s = strings.ReplaceAll(s, "\"", "”")
			return s
		},
	})
	args := map[string]interface{}{
		"res": res,
	}
	switch item := update.Item().(type) {
	// add or remove
	case value.Item:
		if update.Type == value.UpdateTypeAdd {
			tmpl = template.Must(tmpl.Parse(resources.SlackArrayTemplateAdd))
		} else {
			tmpl = template.Must(tmpl.Parse(resources.SlackArrayTemplateRemove))
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
		tmpl = template.Must(tmpl.Parse(resources.SlackArrayTemplateChange))
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

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, args); err != nil {
		return nil, err
	}

	payload := buf.String()
	var payloadValue map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &payloadValue); err != nil {
		ctx.Log.WithField("err", err).Error("Failed to parse payload as json")
		fmt.Println(payload)
		return nil, err
	}
	return payloadValue, nil
}
