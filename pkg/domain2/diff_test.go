package domain2

import (
	"reflect"
	"testing"
)

func TestCompareItemList(t *testing.T) {
	type args struct {
		list1 ItemList
		list2 ItemList
	}
	tests := []struct {
		name string
		args args
		want *Updates
	}{
		{
			args: args{
				list1: ItemList{
					Item{ItemKeyID: "item1", "a": "1", "remove": "2", "change": "3"},
					Item{ItemKeyID: "item2", "c": "1", "d": "3"},
				},
				list2: ItemList{
					Item{ItemKeyID: "item1", "a": "1", "add": "3", "change": "4"},
					Item{ItemKeyID: "item3", "e": "4", "f": "5"},
				},
			},
			want: &Updates{
				Added: []Item{
					Item{ItemKeyID: "item3", "e": "4", "f": "5"},
				},
				Removed: []Item{
					Item{ItemKeyID: "item2", "c": "1", "d": "3"},
				},
				Changed: []ItemChange{
					ItemChange{
						AddedKeys: map[string]string{
							"add": "3",
						},
						RemovedKeys: map[string]string{
							"remove": "2",
						},
						ChangedKeys: map[string]ItemValueChange{
							"change": ItemValueChange{
								Old: "3",
								New: "4",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareItemList(tt.args.list1, tt.args.list2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompareItemList() = %v, want %v", got, tt.want)
			}
		})
	}
}
