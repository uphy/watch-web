package domain

import (
	"strings"
)

type (
	StringValue struct {
		s string
	}
)

func NewStringValue(s string) *StringValue {
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

func (j *StringValue) ItemList() ItemList {
	s := strings.Trim(j.s, " \t\n")
	splitted := strings.Split(s, "\n")
	list := make(ItemList, len(splitted))
	for i, line := range splitted {
		list[i] = Item{line: ""}
	}
	return list
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
