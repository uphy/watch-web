package domain

import (
	"fmt"

	"github.com/uphy/watch-web/pkg/domain2"
)

type (
	JSONObjectValue struct {
		m JSONObject
	}
	JSONArrayValue struct {
		a JSONArray
	}
)

func NewJSONObjectValue(m map[string]interface{}) Value {
	return &JSONObjectValue{m}
}

func NewJSONArrayValue(a []interface{}) Value {
	return &JSONArrayValue{a}
}

func (j *JSONObjectValue) Type() ValueType {
	return ValueTypeJSONObject
}

func (j *JSONObjectValue) JSONObject() JSONObject {
	return j.m
}

func (j *JSONObjectValue) JSONArray() JSONArray {
	return []interface{}{j.m}
}

func (j *JSONObjectValue) ItemList() domain2.ItemList {
	elm := make(domain2.Item)
	for k, v := range j.m {
		elm[k] = fmt.Sprint(v)
	}
	return domain2.ItemList{elm}
}

func (j *JSONObjectValue) String() string {
	return j.m.String()
}

func (j *JSONObjectValue) Interface() interface{} {
	return j.m
}

func (j *JSONObjectValue) Empty() bool {
	return len(j.m) == 0
}

func (j *JSONArrayValue) Type() ValueType {
	return ValueTypeJSONArray
}

func (j *JSONArrayValue) JSONObject() JSONObject {
	return map[string]interface{}{
		"values": j.a,
	}
}

func (j *JSONArrayValue) JSONArray() JSONArray {
	return j.a
}

func (j *JSONArrayValue) ItemList() domain2.ItemList {
	list := make(domain2.ItemList, len(j.a))
	for i, arrayElement := range j.a {
		switch v := arrayElement.(type) {
		case JSONObject:
			elm := make(domain2.Item)
			for k, v := range v {
				elm[k] = fmt.Sprint(v)
			}
			list[i] = elm
		case map[string]interface{}:
			elm := make(domain2.Item)
			for k, v := range v {
				elm[k] = fmt.Sprint(v)
			}
			list[i] = elm
		default:
			list[i] = map[string]string{
				fmt.Sprint(arrayElement): "",
			}
		}
	}
	return list
}

func (j *JSONArrayValue) String() string {
	return j.a.String()
}

func (j *JSONArrayValue) Interface() interface{} {
	return j.a
}

func (j *JSONArrayValue) Empty() bool {
	return len(j.a) == 0
}
