package result

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type (
	DiffResult struct {
		diff []diffmatchpatch.Diff
	}
)

func Diff(v1, v2 string) *DiffResult {
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
	diff := d.DiffCharsToLines(diffs, c)
	return &DiffResult{diff}
}

func DiffJSONArray(jsonArray1, jsonArray2 string) (*DiffResult, error) {
	v1, err := splitJSONArray(jsonArray1)
	if err != nil {
		return nil, err
	}
	v2, err := splitJSONArray(jsonArray2)
	if err != nil {
		return nil, err
	}
	return Diff(v1, v2), nil
}

// splitJSONArray splits the input json array string into elements separated with line break.
func splitJSONArray(jsonArray string) (string, error) {
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

func DiffJSONObject(jsonObject1, jsonObject2 string) (*DiffResult, error) {
	v1, err := splitJSONObject(jsonObject1)
	if err != nil {
		return nil, err
	}
	v2, err := splitJSONObject(jsonObject2)
	if err != nil {
		return nil, err
	}
	return Diff(v1, v2), nil
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

func (d *DiffResult) Changed() bool {
	if len(d.diff) == 0 {
		return false
	}
	for _, d := range d.diff {
		if d.Type == diffmatchpatch.DiffEqual {
			continue
		}
		return true
	}
	return false
}

func (d *DiffResult) String() string {
	w := new(bytes.Buffer)
	for _, diff := range d.diff {
		text := diff.Text
		text = strings.Trim(text, "\r\n")
		texts := strings.Split(text, "\n")
		switch diff.Type {
		case diffmatchpatch.DiffDelete:
			for _, t := range texts {
				fmt.Fprintf(w, "- %s\n", t)
			}
		case diffmatchpatch.DiffEqual:
			for _, t := range texts {
				fmt.Fprintf(w, "  %s\n", t)
			}
		case diffmatchpatch.DiffInsert:
			for _, t := range texts {
				fmt.Fprintf(w, "+ %s\n", t)
			}
		}
	}
	return w.String()
}
