package script

import (
	"reflect"
	"testing"

	"github.com/uphy/watch-web/pkg/domain"
)

func TestJavaScript_Evaluate(t *testing.T) {
	tests := []struct {
		name    string
		script  string
		args    map[string]interface{}
		want    interface{}
		wantErr bool
	}{
		{
			script: "1+1",
			args:   nil,
			want:   2.,
		},
		{
			script: "a+1",
			args: map[string]interface{}{
				"a": 1,
			},
			want: 2.,
		},
		{
			script: "a.replace('a','X')",
			args: map[string]interface{}{
				"a": "abc",
			},
			want: "Xbc",
		},
		{
			script: "a.b",
			args: map[string]interface{}{
				"a": domain.NewJSONObject(map[string]interface{}{
					"b": 1,
				}),
			},
			want: 1,
		},
		{
			script: "a.b *= 2;a.b",
			args: map[string]interface{}{
				"a": domain.NewJSONObject(map[string]interface{}{
					"b": 1,
				}),
			},
			want: 2.,
		},
		{
			script: "a.b *= 2;a",
			args: map[string]interface{}{
				"a": domain.NewJSONObject(map[string]interface{}{
					"b": 1,
				}),
			},
			want: domain.NewJSONObject(map[string]interface{}{
				"b": 2.,
			}),
		},
		{
			script: "a < 2",
			args: map[string]interface{}{
				"a": 1,
			},
			want: true,
		},
	}
	engine := NewJavaScriptEngine()
	for _, tt := range tests {
		script, err := engine.NewScript(tt.script)
		if err != nil {
			t.Errorf("failed to parse script: %v", err)
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := script.Evaluate(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("JavaScript.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JavaScript.Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
