package transformer

import (
	"reflect"
	"testing"

	"github.com/uphy/watch-web/pkg/domain"
	"github.com/uphy/watch-web/pkg/domain/script"
)

func TestScriptTransformer_Transform(t *testing.T) {
	type args struct {
		v      domain.Value
		script string
	}
	tests := []struct {
		name    string
		args    args
		want    domain.Value
		wantErr bool
	}{
		{
			args: args{
				v: domain.NewJSONArray([]interface{}{
					map[string]interface{}{
						"id":    "000",
						"price": 100,
					},
					map[string]interface{}{
						"id":    "001",
						"price": 200,
					},
				}),
				script: `
source.Map(func(v){
	v.description = sprintf("ID: %s (%d yen)", v.id, v.price)
	v.price *= 1.08
	v
}).Filter(func(v){
	v.price < 150
})
`,
			},
			want: domain.NewJSONArray([]interface{}{
				map[string]interface{}{
					"id":          "000",
					"description": "ID: 000 (100 yen)",
					"price":       108.,
				},
			}),
		},
	}
	scriptEngine := script.NewAnkoScriptEngine()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script, err := scriptEngine.NewScript(tt.args.script)
			if err != nil {
				t.Errorf("failed to create script object: %v", err)
				return
			}
			transformer, err := NewScriptTransformer(script)
			if (err != nil) != tt.wantErr {
				t.Errorf("JavaScriptTransformer.Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			ctx := domain.NewDefaultJobContext()
			got, err := transformer.Transform(ctx, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("JavaScriptTransformer.Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JavaScriptTransformer.Transform() = %v, want %v", got, tt.want)
			}
		})
	}
}
