package domain

import (
	"encoding/json"
	"fmt"
	"log"
)

func NewJSONObject(m map[string]interface{}) JSONObject {
	return JSONObject(m)
}

func NewJSONArray(a []interface{}) JSONArray {
	return JSONArray(a)
}

func (j JSONObject) Type() ValueType {
	return ValueTypeJSONObject
}

func (j JSONObject) JSONObject() JSONObject {
	return JSONObject(j)
}

func (j JSONObject) JSONArray() JSONArray {
	return []interface{}{j}
}

func (j JSONObject) ItemList() ItemList {
	elm := make(Item)
	for k, v := range j {
		elm[k] = fmt.Sprint(v)
	}
	return ItemList{elm}
}

func (j JSONObject) Interface() interface{} {
	return j
}

func (j JSONObject) Empty() bool {
	return len(j) == 0
}

func (j JSONObject) String() string {
	b, err := json.Marshal(j)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func (j JSONArray) Type() ValueType {
	return ValueTypeJSONArray
}

func (j JSONArray) JSONObject() JSONObject {
	return map[string]interface{}{
		"values": j,
	}
}

func (j JSONArray) JSONArray() JSONArray {
	return JSONArray(j)
}

func (j JSONArray) ItemList() ItemList {
	list := make(ItemList, len(j))
	for i, arrayElement := range j {
		switch v := arrayElement.(type) {
		case JSONObject:
			elm := make(Item)
			for k, v := range v {
				elm[k] = fmt.Sprint(v)
			}
			list[i] = elm
		case map[string]interface{}:
			elm := make(Item)
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

func (j JSONArray) Filter(filter func(value interface{}) bool) JSONArray {
	filtered := make([]interface{}, 0)
	for _, v := range j {
		if filter(v) {
			filtered = append(filtered, v)
		}
	}
	return NewJSONArray(filtered)
}

func (j JSONArray) Map(mapFunc func(value interface{}) interface{}) JSONArray {
	mapped := make([]interface{}, 0)
	for _, v := range j {
		mapped = append(mapped, mapFunc(v))
	}
	return NewJSONArray(mapped)
}

func (j JSONArray) Interface() interface{} {
	return j
}

func (j JSONArray) Empty() bool {
	return len(j) == 0
}

func (j JSONArray) String() string {
	b, err := json.Marshal(j)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}
