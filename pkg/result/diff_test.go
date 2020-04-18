package result

import (
	"reflect"
	"testing"
)

func Test_diff(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				v1: "aaa\nbbb",
				v2: "aaa\nbbb\nccc",
			},
			want: "  aaa\n  bbb\n+ ccc\n",
		},
		{
			args: args{
				v1: "",
				v2: "aaa",
			},
			want: "+ aaa\n",
		},
		{
			args: args{
				v1: "aaa",
				v2: "",
			},
			want: "- aaa\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Diff(tt.args.v1, tt.args.v2); !reflect.DeepEqual(got.String(), tt.want) {
				t.Errorf("Diff() = '%v', want '%v'", got.String(), tt.want)
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
		want    string
		wantErr bool
	}{
		{
			args: args{
				`["a","b"]`,
				`["a","b","c"]`,
			},
			want: `  "a"
  "b"
+ "c"
`,
			wantErr: false,
		},
		{
			args: args{
				`[{"name":"a"},{"name":"b"}]`,
				`[{"name":"a"},{"name":"b"},{"name":"c"}]`,
			},
			want: `  {"name":"a"}
  {"name":"b"}
+ {"name":"c"}
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DiffJSONArray(tt.args.jsonArray1, tt.args.jsonArray2)
			if (err != nil) != tt.wantErr {
				t.Errorf("DiffJSONArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.String(), tt.want) {
				t.Errorf("DiffJSONArray() = %v, want %v", got.String(), tt.want)
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
		want    string
		wantErr bool
	}{
		{
			args: args{
				`{"name":"a","num":1}`,
				`{"name":"a","num":2,"num2":3}`,
			},
			want: `  {"name":"a"}
- {"num":1}
+ {"num":2}
+ {"num2":3}
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DiffJSONObject(tt.args.jsonObject1, tt.args.jsonObject2)
			if (err != nil) != tt.wantErr {
				t.Errorf("DiffJSONObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.String(), tt.want) {
				t.Errorf("DiffJSONObject() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}
