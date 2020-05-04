package domain

import (
	"reflect"
	"testing"

	"github.com/uphy/watch-web/pkg/domain2"
)

func TestItemList(t *testing.T) {
	tests := []struct {
		name  string
		value Value
		want  domain2.ItemList
	}{
		{
			value: NewJSONObjectValue(map[string]interface{}{
				"a": "A",
				"b": 1,
				"c": true,
			}),
			want: domain2.ItemList{
				domain2.Item{
					"a": "A",
					"b": "1",
					"c": "true",
				},
			},
		},
		{
			value: NewJSONObjectValue(map[string]interface{}{
				"a": "A",
				"b": 1,
				"c": map[string]interface{}{
					"d": 1,
				},
			}),
			want: domain2.ItemList{
				domain2.Item{
					"a": "A",
					"b": "1",
					"c": "map[d:1]",
				},
			},
		},
		{
			value: NewJSONArrayValue([]interface{}{
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
			want: domain2.ItemList{
				domain2.Item{
					"a": "A",
					"b": "1",
					"c": "true",
				},
				domain2.Item{
					"a": "AA",
					"b": "2",
					"c": "false",
				},
			},
		},
		{
			value: NewJSONArrayValue([]interface{}{
				1,
				"a",
			}),
			want: domain2.ItemList{
				domain2.Item{"1": ""},
				domain2.Item{"a": ""},
			},
		},
		{
			value: NewStringValue("aaa"),
			want: domain2.ItemList{
				domain2.Item{"aaa": ""},
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
