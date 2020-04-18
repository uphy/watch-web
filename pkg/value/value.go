package value

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	ValueTypeJSONObject     ValueType = "json_object"
	ValueTypeJSONArray      ValueType = "json_array"
	ValueTypeJSONAutoDetect ValueType = "json"
	ValueTypeString         ValueType = "string"
	ValueTypeAutoDetect     ValueType = "auto"
)

type (
	ValueType string
	Value     interface {
		Type() ValueType
		JSONObject() map[string]interface{}
		JSONArray() []interface{}
		String() string
		Interface() interface{}
		Empty() bool
	}
)

func ConvertAs(s string, valueType ValueType) (Value, error) {
	switch valueType {
	case ValueTypeJSONObject:
		return ParseJSONObject(s)
	case ValueTypeJSONArray:
		return ParseJSONArray(s)
	case ValueTypeString:
		return String(s), nil
	case ValueTypeJSONAutoDetect:
		return ParseJSON(s)
	case ValueTypeAutoDetect:
		return Auto(s), nil
	default:
		return nil, fmt.Errorf("unsupported value type: %s", valueType)
	}
}

func Auto(s string) Value {
	if v, err := ParseJSON(s); err == nil {
		return v
	}
	return String(s)
}

func ParseJSON(s string) (Value, error) {
	if v, err := ParseJSONArray(s); err == nil {
		return v, nil
	} else if v, err := ParseJSONObject(s); err == nil {
		return v, nil
	}
	return nil, errors.New("not a json")
}

func ParseJSONObject(s string) (Value, error) {
	trimmed := strings.Trim(s, " \t\n")
	if strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(s), &m); err == nil {
			return NewJSONObjectValue(m), nil
		}
	}
	return nil, errors.New("not a json object")
}

func ParseJSONArray(s string) (Value, error) {
	trimmed := strings.Trim(s, " \t\n")
	if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
		var a []interface{}
		if err := json.Unmarshal([]byte(s), &a); err == nil {
			return NewJSONArrayValue(a), nil
		}
	}
	return nil, errors.New("not a json array")
}