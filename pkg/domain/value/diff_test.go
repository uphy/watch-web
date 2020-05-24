package value

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {
	updates := Updates{
		*updateChange(
			&ItemChange{
				Item: Item{
					ItemKeyID: "id1",
				},
				AddedKeys: map[string]string{
					"add": "3",
				},
				ChangedKeys: map[string]ItemValueChange{
					"change": {
						Old: "3",
						New: "4",
					},
				},
				RemovedKeys: map[string]string{
					"remove": "2",
				},
			},
		),
		*updateRemove(Item{ItemKeyID: "item2", "c": "1", "d": "3"}),
		*updateAdd(Item{ItemKeyID: "item3", "e": "4", "f": "5"}),
	}
	b, _ := json.MarshalIndent(updates, "", "   ")
	var v Updates
	if err := json.Unmarshal(b, &v); err != nil {
		t.Error("Unmarshal failed:", err)
	}
	if !reflect.DeepEqual(updates, v) {
		t.Error("unmarshal/marshal inconsistent")
	}
}

func TestCompareItemList(t *testing.T) {
	type args struct {
		list1 ItemList
		list2 ItemList
	}
	tests := []struct {
		name string
		args args
		want Updates
	}{
		{
			args: args{
				list1: ItemList{
					Item{ItemKeyID: "item1", "a": "1", "remove": "2", "change": "3", "_ignore1": "1"},
					Item{ItemKeyID: "item2", "c": "1", "d": "3", "_ignored2": "0"},
				},
				list2: ItemList{
					Item{ItemKeyID: "item1", "a": "1", "add": "3", "change": "4"},
					Item{ItemKeyID: "item3", "e": "4", "f": "5"},
				},
			},
			want: Updates{
				*updateChange(
					&ItemChange{
						Item: Item{
							ItemKeyID: "item1",
							"a":       "1",
							"add":     "3",
							"change":  "4",
							"label":   "",
							"link":    "",
							"summary": "",
						},
						AddedKeys: map[string]string{
							"add": "3",
						},
						ChangedKeys: map[string]ItemValueChange{
							"change": {
								Old: "3",
								New: "4",
							},
						},
						RemovedKeys: map[string]string{
							"remove": "2",
						},
					},
				),
				*updateRemove(Item{ItemKeyID: "item2", "c": "1", "d": "3", "label": "", "link": "", "summary": ""}),
				*updateAdd(Item{ItemKeyID: "item3", "e": "4", "f": "5", "label": "", "link": "", "summary": ""}),
			},
		},
		{
			args: args{
				list1: ItemList{
					Item{ItemKeyID: "item1", "line1": ""},
					Item{ItemKeyID: "item2", "line2": ""},
				},
				list2: ItemList{
					Item{ItemKeyID: "item1", "line1": ""},
					Item{ItemKeyID: "item3", "line3": ""},
				},
			},
			want: Updates{
				*updateRemove(Item{ItemKeyID: "item2", "line2": "", "label": "", "link": "", "summary": ""}),
				*updateAdd(Item{ItemKeyID: "item3", "line3": "", "label": "", "link": "", "summary": ""}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareItemList(tt.args.list1, tt.args.list2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompareItemList() = %v, want %v", toJSON(got), toJSON(tt.want))
			}
		})
	}
}

func toJSON(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "   ")
	return string(b)
}
