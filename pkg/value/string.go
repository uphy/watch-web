package value

import "strings"

type (
	StringValue struct {
		s string
	}
)

func String(s string) Value {
	return &StringValue{strings.Trim(s, " \t\n")}
}

func (s *StringValue) Type() ValueType {
	return ValueTypeString
}

func (s *StringValue) JSONObject() JSONObject {
	return map[string]interface{}{
		s.s: "",
	}
}

func (s *StringValue) JSONArray() JSONArray {
	return []interface{}{s.s}
}

func (s *StringValue) String() string {
	return s.s
}

func (s *StringValue) Interface() interface{} {
	return s.s
}

func (s *StringValue) Empty() bool {
	return len(s.s) == 0
}
