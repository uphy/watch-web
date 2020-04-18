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
			if got := diff(tt.args.v1, tt.args.v2); !reflect.DeepEqual(got.String(), tt.want) {
				t.Errorf("diff() = '%v', want '%v'", got.String(), tt.want)
			}
		})
	}
}
