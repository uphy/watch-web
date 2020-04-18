package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/uphy/watch-web/check"
	"github.com/uphy/watch-web/value"
)

const (
	JSONTransformSourceTypeAuto   = "auto"
	JSONTransformSourceTypeArray  = "array"
	JSONTransformSourceTypeObject = "object"
	JSONTransformSourceTypeRaw    = "raw"
)

type (
	FiltersConfig []FilterConfig
	FilterConfig  struct {
		Template      *TemplateString `json:"template,omitempty"`
		DOM           *TemplateString `json:"dom,omitempty"`
		JSONTransform *struct {
			Type      *JSONTransformSourceType  `json:"type,omitempty"`
			Transform map[string]TemplateString `json:"transform,omitempty"`
		} `json:"json_transform,omitempty"`
	}
	JSONTransformSourceType string

	FilterSource struct {
		source  check.Source
		filters []Filter
	}
	Filter interface {
		Filter(v value.Value) (value.Value, error)
	}
	TemplateFilter struct {
		template TemplateString
		ctx      *TemplateContext
	}
	DOMFilter struct {
		selecter string
	}
	JSONTransformFilter struct {
		sourceType JSONTransformSourceType
		transform  map[string]TemplateString
		ctx        *TemplateContext
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
	if f.JSONTransform != nil {
		var t JSONTransformSourceType
		if f.JSONTransform.Type != nil {
			t = *f.JSONTransform.Type
		} else {
			t = JSONTransformSourceTypeAuto
		}
		return &JSONTransformFilter{t, f.JSONTransform.Transform, ctx}, nil
	}
	return nil, errors.New("no filters defined")
}

func (f *FilterSource) Fetch() (value.Value, error) {
	v, err := f.source.Fetch()
	if err != nil {
		return nil, err
	}
	for _, filter := range f.filters {
		filtered, err := filter.Filter(v)
		if err != nil {
			return nil, err
		}
		v = filtered
	}
	return v, nil
}

func (f *FilterSource) String() string {
	return fmt.Sprintf("Filter[source=%v, filters=%v]", f.source, f.filters)
}

func (t *TemplateFilter) Filter(v value.Value) (value.Value, error) {
	t.ctx.PushScope()
	defer t.ctx.PopScope()
	t.ctx.Set("source", v)
	evaluated, err := t.template.Evaluate(t.ctx)
	if err != nil {
		return nil, err
	}
	return value.Auto(evaluated), nil
}

func (t *TemplateFilter) String() string {
	return fmt.Sprintf("Template[template=%v]", t.template)
}

func (t *DOMFilter) Filter(v value.Value) (value.Value, error) {
	return parseDOM(v.String(), t.selecter)
}

func parseDOM(html string, selector string) (value.Value, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	selection := doc.Find(selector)

	nodes := make([]interface{}, 0)
	for _, node := range selection.Nodes {
		var v = make(map[string]interface{})
		v["data"] = node.Data
		for _, a := range node.Attr {
			v[a.Key] = a.Val
		}
		nodes = append(nodes, v)
	}
	result["text"] = selection.Text()
	selectedHTML, _ := selection.Html()
	result["html"] = selectedHTML
	result["nodes"] = nodes
	return value.NewJSONObjectValue(result), nil
}

func (t *DOMFilter) String() string {
	return fmt.Sprintf("DOM[selector=%v]", t.selecter)
}

func (t *JSONTransformFilter) Filter(v value.Value) (value.Value, error) {
	// auto detect source type
	var sourceType = t.sourceType
	if sourceType == JSONTransformSourceTypeAuto {
		s := v.String()
		if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
			sourceType = JSONTransformSourceTypeArray
		} else if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
			sourceType = JSONTransformSourceTypeObject
		} else {
			sourceType = JSONTransformSourceTypeRaw
		}
	}
	// parse source based on type
	var source []interface{}
	hasMultiElements := false
	switch sourceType {
	case JSONTransformSourceTypeArray:
		source = v.JSONArray()
		hasMultiElements = true
	case JSONTransformSourceTypeObject:
		source = []interface{}{v.JSONObject()}
	case JSONTransformSourceTypeRaw:
		source = []interface{}{v.Interface()}
	default:
		return nil, fmt.Errorf("unsupported transform source type: %v", sourceType)
	}
	// transform
	var transformed []interface{}
	for _, src := range source {
		t.ctx.PushScope()
		t.ctx.Set("source", src)
		elm := make(map[string]string)
		for k, tmpl := range t.transform {
			evaluated, err := tmpl.Evaluate(t.ctx)
			if err != nil {
				t.ctx.PopScope()
				return nil, err
			}
			elm[k] = evaluated
		}
		transformed = append(transformed, elm)
		t.ctx.PopScope()
	}
	// return
	if hasMultiElements {
		return value.NewJSONArrayValue(transformed), nil
	}
	return value.NewJSONObjectValue(transformed[0].(map[string]interface{})), nil
}
