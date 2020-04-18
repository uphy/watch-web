package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/PuerkitoBio/goquery"
	"github.com/uphy/watch-web/pkg/value"
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
		v, err := ParseDOM(html, selector)
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
}

func (t TemplateString) Evaluate(ctx *TemplateContext) (string, error) {
	var buf = new(bytes.Buffer)
	if err := template.Must(template.New("template-string").Funcs(funcs).Parse(string(t))).Execute(buf, ctx.Get()); err != nil {
		return "", fmt.Errorf("cannot evaluate %s: %w", t, err)
	}
	return buf.String(), nil
}

func ParseDOM(html string, selector string) (value.Value, error) {
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
	return value.NewJSONObjectValue(result), nil
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
