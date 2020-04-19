package domain

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

func (j *JSONArrayValue) String() string {
	return j.a.String()
}

func (j *JSONArrayValue) Interface() interface{} {
	return j.a
}

func (j *JSONArrayValue) Empty() bool {
	return len(j.a) == 0
}
