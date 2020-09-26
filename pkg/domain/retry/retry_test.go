package retry

import (
	"errors"
	"testing"
)

func TestRetrier_Run(t *testing.T) {
	type fields struct {
		initialInterval     float64
		multiplier          float64
		randomizationFactor float64
		retry               int
	}
	type args struct {
		task RetryableFunc
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "no retry",
			fields: fields{
				initialInterval:     1,
				multiplier:          2,
				randomizationFactor: 1,
				retry:               1,
			},
			args: args{
				task: func(ctx *RetryContext) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "retry - success",
			fields: fields{
				initialInterval:     1,
				multiplier:          2,
				randomizationFactor: 1,
				retry:               1,
			},
			args: args{
				task: func(ctx *RetryContext) error {
					switch ctx.Retried {
					case 0:
						return errors.New("error")
					default:
						return nil
					}
				},
			},
			wantErr: false,
		},
		{
			name: "retry - fail",
			fields: fields{
				initialInterval:     1,
				multiplier:          2,
				randomizationFactor: 1,
				retry:               1,
			},
			args: args{
				task: func(ctx *RetryContext) error {
					return errors.New("error")
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Retrier{
				retry:               tt.fields.retry,
				initialInterval:     tt.fields.initialInterval,
				randomizationFactor: tt.fields.randomizationFactor,
				multiplier:          tt.fields.multiplier,
			}
			if err := r.Run(tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("Retrier.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
