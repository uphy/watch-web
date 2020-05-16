package value

import (
	"reflect"
	"testing"
)

func TestItemList(t *testing.T) {
	tests := []struct {
		name  string
		value Value
		want  ItemList
	}{
		{
			value: NewJSONObject(map[string]interface{}{
				"a": "A",
				"b": 1,
				"c": true,
			}),
			want: ItemList{
				Item{
					"a": "A",
					"b": "1",
					"c": "true",
				},
			},
		},
		{
			value: NewJSONObject(map[string]interface{}{
				"a": "A",
				"b": 1,
				"c": map[string]interface{}{
					"d": 1,
				},
			}),
			want: ItemList{
				Item{
					"a": "A",
					"b": "1",
					"c": "map[d:1]",
				},
			},
		},
		{
			value: NewJSONArray([]interface{}{
				map[string]interface{}{
					"a": "A",
					"b": 1,
					"c": true,
				},
				map[string]interface{}{
					"a": "AA",
					"b": 2,
					"c": false,
				},
			}),
			want: ItemList{
				Item{
					"a": "A",
					"b": "1",
					"c": "true",
				},
				Item{
					"a": "AA",
					"b": "2",
					"c": "false",
				},
			},
		},
		{
			value: NewJSONArray([]interface{}{
				1,
				"a",
			}),
			want: ItemList{
				Item{"1": ""},
				Item{"a": ""},
			},
		},
		{
			value: NewStringValue("aaa"),
			want: ItemList{
				Item{"aaa": ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := tt.value
			if got := j.ItemList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONObjectValue.ItemList() = %v, want %v", got, tt.want)
			}
		})
	}
}
