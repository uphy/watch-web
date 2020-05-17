package action

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

	"github.com/uphy/watch-web/pkg/domain/value"

	"github.com/uphy/watch-web/pkg/domain"
	"github.com/uphy/watch-web/pkg/resources"

	"gopkg.in/yaml.v2"
)

const (
	DiffSplitSize = 5
)

type (
	SlackAction struct {
		URL   string
		Debug bool
	}
)

func NewSlackAction(webhookURL string, debug bool) *SlackAction {
	return &SlackAction{webhookURL, debug}
}

func (s *SlackAction) Run(ctx *domain.JobContext, res *domain.Result) error {
	splittedUpdates := s.splitUpdates(res.Diff())
	for _, updates := range splittedUpdates {
		if !updates.Changes() {
			continue
		}
		payloadValue, err := s.run(ctx, res, updates)
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
			if resp.StatusCode != 200 {
				b, _ := ioutil.ReadAll(resp.Body)
				return fmt.Errorf("invalid status code: status=%d, body=%s", resp.StatusCode, string(b))
			}
		}
	}
	return nil
}

func (s *SlackAction) splitUpdates(updates value.Updates) []value.Updates {
	splitted := make([]value.Updates, 0)
	tmpUpdates := make(value.Updates, 0)
	for _, update := range updates {
		tmpUpdates = append(tmpUpdates, update)
		if len(tmpUpdates) >= DiffSplitSize {
			splitted = append(splitted, tmpUpdates)
			tmpUpdates = make(value.Updates, 0)
		}
	}
	if len(tmpUpdates) >= 0 {
		splitted = append(splitted, tmpUpdates)
	}
	return splitted
}

func (s *SlackAction) run(ctx *domain.JobContext, res *domain.Result, updates value.Updates) (interface{}, error) {
	v := map[string]interface{}{
		"res":     res,
		"updates": updates,
	}
	tmpl := template.Must(template.New("slack-template").Funcs(template.FuncMap{
		"toArrayExcludes": func(v map[string]string, excludes ...string) []map[string]interface{} {
			excludesSet := make(map[string]struct{})
			for _, exclude := range excludes {
				excludesSet[exclude] = struct{}{}
			}
			var result = make([]map[string]interface{}, 0)
			for k, v := range v {
				if _, isExclude := excludesSet[k]; isExclude {
					continue
				}
				result = append(result, map[string]interface{}{
					"key":   k,
					"value": v,
				})
			}
			return result
		},
		"flatChanges": func(changes *value.ItemChange) []map[string]interface{} {
			var result = make([]map[string]interface{}, 0)
			for k, v := range changes.AddedKeys {
				if k == value.ItemKeyID {
					continue
				}
				result = append(result, map[string]interface{}{
					"key":   k,
					"value": v,
					"type":  "add",
				})
			}
			for k, v := range changes.RemovedKeys {
				if k == value.ItemKeyID {
					continue
				}
				result = append(result, map[string]interface{}{
					"key":   k,
					"value": v,
					"type":  "remove",
				})
			}
			for k, v := range changes.ChangedKeys {
				if k == value.ItemKeyID {
					continue
				}
				result = append(result, map[string]interface{}{
					"key":   k,
					"value": v,
					"type":  "change",
				})
			}
			return result
		},
		"escape": func(s string) string {
			s = strings.ReplaceAll(s, "\n", "\\n")
			return s
		},
	}).Parse(resources.SlackArrayTemplate))
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, v); err != nil {
		return nil, err
	}

	payload := buf.String()
	var payloadValue interface{}
	if err := json.Unmarshal([]byte(payload), &payloadValue); err != nil {
		ctx.Log.WithField("err", err).Error("Failed to parse payload as json")
		fmt.Println(payload)
		return nil, err
	}
	return payloadValue, nil
}
