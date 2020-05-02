package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	ChangeTypeInsert ChangeType = "insert"
	ChangeTypeDelete ChangeType = "delete"
	ChangeTypeEqual  ChangeType = "equal"
)

type (
	DiffResult interface {
		Changed() bool
		String() string
	}

	ChangeType       string
	StringDiffResult []Line
	Line             struct {
		Text string
		Type ChangeType
	}
	JSONObjectDiffResult []JSONField
	JSONField            struct {
		Name  string
		Value interface{}
		Type  ChangeType
	}
	JSONArrayDiffResult []JSONArrayElement
	JSONArrayElement    struct {
		Object JSONObject
		Type   ChangeType
	}
)

func DiffString(v1, v2 string) StringDiffResult {
	d := diff(v1, v2)
	lines := make([]Line, 0)
	for _, l := range d {
		var t ChangeType
		switch l.Type {
		case diffmatchpatch.DiffInsert:
			t = ChangeTypeInsert
		case diffmatchpatch.DiffDelete:
			t = ChangeTypeDelete
		case diffmatchpatch.DiffEqual:
			t = ChangeTypeEqual
		default:
			log.Fatal("unexpected type: ", l.Type)
		}
		texts := strings.Split(strings.TrimRight(l.Text, "\n"), "\n")
		for _, text := range texts {
			lines = append(lines, Line{text, t})
		}
	}
	return lines
}

func diff(v1, v2 string) []diffmatchpatch.Diff {
	v1 = strings.Trim(v1, " \t\n")
	if len(v1) > 0 {
		v1 = v1 + "\n"
	}
	v2 = strings.Trim(v2, " \t\n")
	if len(v2) > 0 {
		v2 = v2 + "\n"
	}

	d := diffmatchpatch.New()
	a, b, c := d.DiffLinesToChars(v1, v2)
	diffs := d.DiffMain(a, b, false)
	return d.DiffCharsToLines(diffs, c)
}

func DiffJSONArray(jsonArray1, jsonArray2 string) (JSONArrayDiffResult, error) {
	v1, err := splitJSONArray(jsonArray1)
	if err != nil {
		return nil, err
	}
	v2, err := splitJSONArray(jsonArray2)
	if err != nil {
		return nil, err
	}
	d := DiffString(v1, v2)
	elements := make([]JSONArrayElement, len(d))
	for i, line := range d {
		var v JSONObject
		if err := json.Unmarshal([]byte(line.Text), &v); err != nil {
			return nil, err
		}
		elements[i] = JSONArrayElement{v, line.Type}
	}
	return elements, nil
}

func DiffJSONObject(jsonObject1, jsonObject2 string) (JSONObjectDiffResult, error) {
	v1, err := splitJSONObject(jsonObject1)
	if err != nil {
		return nil, err
	}
	v2, err := splitJSONObject(jsonObject2)
	if err != nil {
		return nil, err
	}
	d := DiffString(v1, v2)
	fields := make([]JSONField, len(d))
	for i, line := range d {
		var f JSONObject
		if err := json.Unmarshal([]byte(line.Text), &f); err != nil {
			return nil, err
		}
		if len(f) != 1 {
			return nil, fmt.Errorf("expected a field but %v", f)
		}
		for k, v := range f {
			fields[i] = JSONField{k, v, line.Type}
			break
		}
	}
	sort.Slice(fields, func(i, j int) bool {
		return strings.Compare(fields[i].Name, fields[j].Name) < 0
	})
	return fields, nil
}

// splitJSONArray splits the input json array string into elements separated with line break.
func splitJSONArray(jsonArray string) (string, error) {
	if jsonArray == "" {
		// Replace empty string with empty array
		// because initial state is empty string "".
		jsonArray = "[]"
	}
	var v []interface{}
	if err := json.Unmarshal([]byte(jsonArray), &v); err != nil {
		return "", err
	}
	var s []string
	for _, elm := range v {
		b, err := json.Marshal(elm)
		if err != nil {
			return "", err
		}
		s = append(s, string(b))
	}
	return strings.Join(s, "\n"), nil
}

// splitJSONObject splits the input json object string into single field objects separated with line break
func splitJSONObject(jsonObject string) (string, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(jsonObject), &obj); err != nil {
		return "", err
	}
	var s []string
	for k, v := range obj {
		m := map[string]interface{}{
			k: v,
		}
		b, err := json.Marshal(m)
		if err != nil {
			return "", err
		}
		s = append(s, string(b))
	}
	return strings.Join(s, "\n"), nil
}

func (r StringDiffResult) Changed() bool {
	for _, l := range r {
		if l.Type != ChangeTypeEqual {
			return true
		}
	}
	return false
}

func (r StringDiffResult) String() string {
	w := new(bytes.Buffer)
	for _, l := range r {
		fmt.Fprintln(w, diffLineToString(l.Text, l.Type))
	}
	return w.String()
}

func diffLineToString(line string, t ChangeType) string {
	switch t {
	case ChangeTypeInsert:
		return "+ " + line
	case ChangeTypeDelete:
		return "- " + line
	default:
		return "  " + line
	}
}

func (r JSONObjectDiffResult) Changed() bool {
	for _, l := range r {
		if l.Type != ChangeTypeEqual {
			return true
		}
	}
	return false
}

func (r JSONObjectDiffResult) String() string {
	w := new(bytes.Buffer)
	for _, l := range r {
		fmt.Fprintln(w, diffLineToString(fmt.Sprintf("%s = %v", l.Name, l.Value), l.Type))
	}
	return w.String()
}

func (r JSONArrayDiffResult) Changed() bool {
	for _, l := range r {
		if l.Type != ChangeTypeEqual {
			return true
		}
	}
	return false
}

func (r JSONArrayDiffResult) String() string {
	w := new(bytes.Buffer)
	for _, l := range r {
		fmt.Fprintln(w, diffLineToString(fmt.Sprintf("%v", l.Object), l.Type))
	}
	return w.String()
}
