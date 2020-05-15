package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type (
	TemplateString string
)

var funcs = map[string]interface{}{
	"env": func(name string) string {
		return os.Getenv(name)
	},
	"default": func(defaultValue string, value string) string {
		if value == "" {
			return defaultValue
		}
		return value
	},
	"sliceOf": func(values ...string) []string {
		return values
	},
	"json": func(jsonString string) (interface{}, error) {
		if strings.Trim(jsonString, " ") == "" {
			return make(map[string]interface{}), nil
		}
		var v interface{}
		if err := json.Unmarshal([]byte(jsonString), &v); err != nil {
			return nil, err
		}
		return v, nil
	},
	"dom": func(selector, html string) (interface{}, error) {
		// this returns parsed DOM as map(json object value)
		v, err := SelectDOM(html, selector)
		if err != nil {
			return nil, err
		}
		return v.Interface(), nil
	},
	"trim": func(s string) string {
		return strings.Trim(s, " 　\t\r\n")
	},
	"jsonFormat": func(jsonString string) (string, error) {
		var v interface{}
		if err := json.Unmarshal([]byte(jsonString), &v); err != nil {
			return "", err
		}
		b, _ := json.MarshalIndent(v, "", "   ")
		return string(b), nil
	},
	"truncate": func(length int, s string) string {
		runes := []rune(s)
		if len(runes) > length {
			return string(runes[0:length]) + "..."
		}
		return s
	},
	"match": func(pattern string, s string) (string, error) {
		r, err := regexp.Compile(pattern)
		if err != nil {
			return "", err
		}
		matches := r.FindStringSubmatch(s)
		switch len(matches) {
		case 0, 1:
			return "<no match>", nil
		case 2:
			return matches[1], nil
		default:
			return fmt.Sprintf("<multiple matches:%v>", matches), nil
		}
	},
	"replace": func(old, new, original string) string {
		return strings.ReplaceAll(original, old, new)
	},
	"contains": func(substring, s string) bool {
		return strings.Contains(s, substring)
	},
	"formatEpochMillis": func(epochMillis float64) string {
		n := int64(epochMillis) * 1000000
		return time.Unix(0, n).Format("2006/01/02 15:04")
	},
}

func (t TemplateString) Evaluate(ctx *TemplateContext) (string, error) {
	var buf = new(bytes.Buffer)
	tmpl, err := template.New("template-string").Funcs(funcs).Parse(string(t))
	if err != nil {
		return "", fmt.Errorf("cannot parse template: template=%v, error=%w", t, err)
	}
	if err := tmpl.Execute(buf, ctx.Get()); err != nil {
		return "", fmt.Errorf("cannot evaluate %s: %w", t, err)
	}
	return buf.String(), nil
}

func SelectDOM(html string, selector string) (Value, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	selection := doc.Find(selector)

	nodes := make([]interface{}, 0)
	for _, node := range selection.Nodes {
		var v = nodeToMap(node)
		nodes = append(nodes, v)
	}
	result["text"] = selection.Text()
	selectedHTML, _ := selection.Html()
	result["html"] = selectedHTML
	result["nodes"] = nodes
	return NewJSONObject(result), nil
}

func nodeToMap(n *html.Node) map[string]interface{} {
	var v = make(map[string]interface{})
	v["data"] = n.Data
	for _, a := range n.Attr {
		v[a.Key] = a.Val
	}
	children := make([]interface{}, 0)
	if n.FirstChild != nil {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			children = append(children, nodeToMap(c))
		}
	}
	v["children"] = children
	return v
}
