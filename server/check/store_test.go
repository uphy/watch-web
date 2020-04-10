package check

import "testing"

func Test_parseRedisToGoURL(t *testing.T) {
	type args struct {
		redisToGo string
	}
	tests := []struct {
		name         string
		args         args
		wantAddr     string
		wantPassword string
	}{
		{
			name: "",
			args: args{
				redisToGo: "redis://redistogo:xyz@foo.com:10992/",
			},
			wantAddr:     "foo.com:10992",
			wantPassword: "xyz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddr, gotPassword := parseRedisToGoURL(tt.args.redisToGo)
			if gotAddr != tt.wantAddr {
				t.Errorf("parseRedisToGoURL() gotAddr = %v, want %v", gotAddr, tt.wantAddr)
			}
			if gotPassword != tt.wantPassword {
				t.Errorf("parseRedisToGoURL() gotPassword = %v, want %v", gotPassword, tt.wantPassword)
			}
		})
	}
}
