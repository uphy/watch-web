package domain

import (
	"strings"

	"github.com/uphy/watch-web/pkg/domain2"
)

type (
	StringValue struct {
		s string
	}
)

func NewStringValue(s string) Value {
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

func (j *StringValue) ItemList() domain2.ItemList {
	s := strings.Trim(j.s, " \t\n")
	splitted := strings.Split(s, "\n")
	list := make(domain2.ItemList, len(splitted))
	for i, line := range splitted {
		list[i] = domain2.Item{line: ""}
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
