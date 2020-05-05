package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/uphy/watch-web/pkg/domain2"
	"github.com/uphy/watch-web/pkg/resources"

	"github.com/uphy/watch-web/pkg/domain"
	"gopkg.in/yaml.v2"
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
	payloadValue, err := s.run(ctx, res)
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
	return nil
}

func (s *SlackAction) run(ctx *domain.JobContext, res *domain.Result) (interface{}, error) {
	updates := res.Diff()
	if !updates.Changes() {
		return nil, nil
	}
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
		"flatChanges": func(changes *domain2.ItemChange) []map[string]interface{} {
			var result = make([]map[string]interface{}, 0)
			for k, v := range changes.AddedKeys {
				if k == domain2.ItemKeyID {
					continue
				}
				result = append(result, map[string]interface{}{
					"key":   k,
					"value": v,
					"type":  "add",
				})
			}
			for k, v := range changes.RemovedKeys {
				if k == domain2.ItemKeyID {
					continue
				}
				result = append(result, map[string]interface{}{
					"key":   k,
					"value": v,
					"type":  "remove",
				})
			}
			for k, v := range changes.ChangedKeys {
				if k == domain2.ItemKeyID {
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
