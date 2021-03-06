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
		JSONObject() JSONObject
		JSONArray() JSONArray
		ItemList() ItemList
		String() string
		Interface() interface{}
		Empty() bool
	}
	JSONObject map[string]interface{}
	JSONArray  []interface{}
)

func ConvertInterfaceAs(i interface{}, valueType ValueType) (Value, error) {
	if s, ok := i.(string); ok {
		return ConvertAs(s, valueType)
	}
	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	return ConvertAs(string(b), valueType)
}

func ConvertAs(s string, valueType ValueType) (Value, error) {
	switch valueType {
	case ValueTypeJSONObject:
		return ParseJSONObject(s)
	case ValueTypeJSONArray:
		return ParseJSONArray(s)
	case ValueTypeString:
		return NewStringValue(s), nil
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
	return NewStringValue(s)
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
	if len(trimmed) == 0 {
		return NewJSONObject(make(JSONObject)), nil
	}
	if strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(s), &m); err == nil {
			return NewJSONObject(m), nil
		}
	}
	return nil, errors.New("not a json object")
}

func ParseJSONArray(s string) (Value, error) {
	trimmed := strings.Trim(s, " \t\n")
	if len(trimmed) == 0 {
		return NewJSONArray(make(JSONArray, 0)), nil
	}
	if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
		var a []interface{}
		if err := json.Unmarshal([]byte(s), &a); err == nil {
			return NewJSONArray(a), nil
		}
	}
	return nil, errors.New("not a json array")
}
