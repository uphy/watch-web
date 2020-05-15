package domain

import (
	"strings"
)

type (
	StringValue string
)

func NewStringValue(s string) StringValue {
	return StringValue(strings.Trim(s, " \t\n"))
}

func (s StringValue) Type() ValueType {
	return ValueTypeString
}

func (s StringValue) JSONObject() JSONObject {
	return map[string]interface{}{
		string(s): "",
	}
}

func (s StringValue) JSONArray() JSONArray {
	return []interface{}{string(s)}
}

func (j StringValue) ItemList() ItemList {
	s := strings.Trim(string(j), " \t\n")
	splitted := strings.Split(s, "\n")
	list := make(ItemList, len(splitted))
	for i, line := range splitted {
		list[i] = Item{line: ""}
	}
	return list
}

func (s StringValue) String() string {
	return string(s)
}

func (s StringValue) Interface() interface{} {
	return string(s)
}

func (s StringValue) Empty() bool {
	return len(s) == 0
}
