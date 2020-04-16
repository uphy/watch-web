package value

import (
	"encoding/json"
	"strings"
)

type (
	Value interface {
		JSONObject() map[string]interface{}
		JSONArray() []interface{}
		String() string
		Interface() interface{}
		Empty() bool
	}
)

func Auto(s string) Value {
	trimmed := strings.Trim(s, " \t\n")
	if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
		var a []interface{}
		if err := json.Unmarshal([]byte(s), &a); err == nil {
			return NewJSONArrayValue(a)
		}
	} else if strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(s), &m); err == nil {
			return NewJSONObjectValue(m)
		}
	}
	return String(s)
}
