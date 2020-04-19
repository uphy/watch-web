package result

import (
	"reflect"
	"testing"

	"github.com/uphy/watch-web/pkg/value"
)

func TestDiffString(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name string
		args args
		want StringDiffResult
	}{
		{
			args: args{
				v1: "aaa\nbbb",
				v2: "aaa\nbbb\nccc",
			},
			want: StringDiffResult{
				Line{
					Text: "aaa",
					Type: ChangeTypeEqual,
				},
				Line{
					Text: "bbb",
					Type: ChangeTypeEqual,
				},
				Line{
					Text: "ccc",
					Type: ChangeTypeInsert,
				},
			},
		},
		{
			args: args{
				v1: "",
				v2: "aaa",
			},
			want: StringDiffResult{
				Line{
					Text: "aaa",
					Type: ChangeTypeInsert,
				},
			},
		},
		{
			args: args{
				v1: "aaa",
				v2: "",
			},
			want: StringDiffResult{
				Line{
					Text: "aaa",
					Type: ChangeTypeDelete,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DiffString(tt.args.v1, tt.args.v2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DiffString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiffJSONArray(t *testing.T) {
	type args struct {
		jsonArray1 string
		jsonArray2 string
	}
	tests := []struct {
		name    string
		args    args
		want    JSONArrayDiffResult
		wantErr bool
	}{
		{
			args: args{
				jsonArray1: `[{"name":"a"},{"name":"b"}]`,
				jsonArray2: `[{"name":"a"},{"name":"b"},{"name":"c"}]`,
			},
			want: JSONArrayDiffResult{
				JSONArrayElement{
					Object: value.JSONObject{
						"name": "a",
					},
					Type: ChangeTypeEqual,
				},
				JSONArrayElement{
					Object: value.JSONObject{
						"name": "b",
					},
					Type: ChangeTypeEqual,
				},
				JSONArrayElement{
					Object: value.JSONObject{
						"name": "c",
					},
					Type: ChangeTypeInsert,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DiffJSONArray(tt.args.jsonArray1, tt.args.jsonArray2)
			if (err != nil) != tt.wantErr {
				t.Errorf("DiffJSONArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DiffJSONArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiffJSONObject(t *testing.T) {
	type args struct {
		jsonObject1 string
		jsonObject2 string
	}
	tests := []struct {
		name    string
		args    args
		want    JSONObjectDiffResult
		wantErr bool
	}{
		{
			args: args{
				jsonObject1: `{"name":"a","num":1}`,
				jsonObject2: `{"name":"a","num":2,"num2":3}`,
			},
			want: JSONObjectDiffResult{
				JSONField{
					Name:  "name",
					Value: "a",
					Type:  ChangeTypeEqual,
				},
				JSONField{
					Name:  "num",
					Value: float64(1),
					Type:  ChangeTypeDelete,
				},
				JSONField{
					Name:  "num",
					Value: float64(2),
					Type:  ChangeTypeInsert,
				},
				JSONField{
					Name:  "num2",
					Value: float64(3),
					Type:  ChangeTypeInsert,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DiffJSONObject(tt.args.jsonObject1, tt.args.jsonObject2)
			if (err != nil) != tt.wantErr {
				t.Errorf("DiffJSONObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DiffJSONObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
