package value

import (
	"encoding/json"
	"log"
)

type (
	JSONObjectValue struct {
		m map[string]interface{}
	}
	JSONArrayValue struct {
		a []interface{}
	}
)

func NewJSONObjectValue(m map[string]interface{}) Value {
	return &JSONObjectValue{m}
}

func NewJSONArrayValue(a []interface{}) Value {
	return &JSONArrayValue{a}
}

func (j *JSONObjectValue) JSONObject() map[string]interface{} {
	return j.m
}

func (j *JSONObjectValue) JSONArray() []interface{} {
	return []interface{}{j.m}
}

func (j *JSONObjectValue) String() string {
	b, err := json.Marshal(j.m)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func (j *JSONObjectValue) Interface() interface{} {
	return j.m
}

func (j *JSONObjectValue) Empty() bool {
	return len(j.m) == 0
}

func (j *JSONArrayValue) JSONObject() map[string]interface{} {
	return map[string]interface{}{
		"values": j.a,
	}
}

func (j *JSONArrayValue) JSONArray() []interface{} {
	return j.a
}

func (j *JSONArrayValue) String() string {
	b, err := json.Marshal(j.a)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func (j *JSONArrayValue) Interface() interface{} {
	return j.a
}

func (j *JSONArrayValue) Empty() bool {
	return len(j.a) == 0
}
