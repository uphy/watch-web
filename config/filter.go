package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/uphy/watch-web/check"
	"golang.org/x/net/html"
)

type (
	FiltersConfig []FilterConfig
	FilterConfig  struct {
		Template *TemplateString `json:"template,omitempty"`
		DOM      *TemplateString `json:"dom,omitempty"`
	}
	FilterSource struct {
		source  check.Source
		filters []Filter
	}
	Filter interface {
		Filter(s string) (string, error)
	}
	TemplateFilter struct {
		template TemplateString
		ctx      *TemplateContext
	}
	DOMFilter struct {
		selecter string
	}
)

func (f FiltersConfig) Filters(ctx *TemplateContext, source check.Source) (check.Source, error) {
	if len(f) == 0 {
		return source, nil
	}
	filters := make([]Filter, 0)
	for _, filterConfig := range f {
		filter, err := filterConfig.Filter(ctx)
		if err != nil {
			return nil, err
		}
		filters = append(filters, filter)
	}
	return &FilterSource{source, filters}, nil
}

func (f *FilterConfig) Filter(ctx *TemplateContext) (Filter, error) {
	if f.Template != nil {
		return &TemplateFilter{*f.Template, ctx}, nil
	}
	if f.DOM != nil {
		selector, err := f.DOM.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		return &DOMFilter{selector}, nil
	}
	return nil, errors.New("no filters defined")
}

func (f *FilterSource) Fetch() (string, error) {
	s, err := f.source.Fetch()
	if err != nil {
		return "", err
	}
	for _, filter := range f.filters {
		filtered, err := filter.Filter(s)
		if err != nil {
			return "", err
		}
		s = filtered
	}
	return s, nil
}

func (f *FilterSource) String() string {
	return fmt.Sprintf("Filter[source=%v, filters=%v]", f.source, f.filters)
}

func (t *TemplateFilter) Filter(s string) (string, error) {
	t.ctx.PushScope()
	defer t.ctx.PopScope()
	t.ctx.Set("source", s)
	evaluated, err := t.template.Evaluate(t.ctx)
	if err != nil {
		return "", err
	}
	return evaluated, nil
}

func (t *TemplateFilter) String() string {
	return fmt.Sprintf("Template[template=%v]", t.template)
}

func (t *DOMFilter) Filter(s string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
	if err != nil {
		return "", err
	}
	result := make(map[string]interface{})
	selection := doc.Find(t.selecter)

	nodes := make([]interface{}, 0)
	for _, node := range selection.Nodes {
		var v = nodeToMap(node)
		nodes = append(nodes, v)
	}
	result["text"] = selection.Text()
	result["nodes"] = nodes
	b, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(b), nil
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

func (t *DOMFilter) String() string {
	return fmt.Sprintf("DOM[selector=%v]", t.selecter)
}
