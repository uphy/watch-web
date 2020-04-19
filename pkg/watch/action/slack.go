package actions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/uphy/watch-web/pkg/domain"
	"github.com/uphy/watch-web/pkg/resources"
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
	changes, err := res.Diff()
	if err != nil {
		return err
	}
	if !changes.Changed() {
		return nil
	}
	var template string
	switch changes.(type) {
	case domain.JSONArrayDiffResult:
		template = resources.SlackArrayTemplate
	case domain.JSONObjectDiffResult:
		template = resources.SlackTemplate
	default:
		template = resources.SlackTemplate
	}
	return s.run(ctx, map[string]interface{}{
		"res":     res,
		"changes": changes,
	}, template)
}

func (s *SlackAction) run(ctx *domain.JobContext, v interface{}, templateString string) error {
	tmpl := template.Must(template.New("slack-template").Funcs(template.FuncMap{
		"toArrayExcludes": func(v map[string]interface{}, excludes ...string) []map[string]interface{} {
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
	}).Parse(templateString))
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, v); err != nil {
		return err
	}
	if s.Debug {
		payload := buf.String()
		ctx.Log.Info("Slack action is debug mode.  No notification.")
		var v interface{}
		if err := json.Unmarshal([]byte(payload), &v); err != nil {
			ctx.Log.WithField("err", err).Error("Failed to parse payload as json")
			fmt.Println(payload)
		}
		b, _ := json.MarshalIndent(v, "", "   ")
		fmt.Println("Payload: ")
		fmt.Println(string(b))
	} else {
		resp, err := http.Post(s.URL, "application/json", bytes.NewReader(buf.Bytes()))
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			return errors.New("invalid status code: " + strconv.Itoa(resp.StatusCode))
		}
	}
	return nil
}
