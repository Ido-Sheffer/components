package http

import (
	"github.com/kubemq-hub/components/config"
	"reflect"
	"testing"
)

func Test_parseOptions(t *testing.T) {
	type args struct {
		cfg config.Metadata
	}
	tests := []struct {
		name    string
		args    args
		want    options
		wantErr bool
	}{
		{
			name: "valid options",
			args: args{
				cfg: config.Metadata{
					Name: "target.http",
					Kind: "target.http",
					Properties: map[string]string{
						"uri":       "http://localhost",
						"auth_type": "basic",
						"username":  "some-user-name",
						"password":  "some-password",
						"token":     "some-token",
						"headers":   `{"header":"value"}`,
					},
				},
			},
			want: options{
				uri:      "http://localhost",
				authType: "basic",
				username: "some-user-name",
				password: "some-password",
				token:    "some-token",
				headers: map[string]string{
					"header": "value",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid options - bad headers format",
			args: args{
				cfg: config.Metadata{
					Name: "target.http",
					Kind: "target.http",
					Properties: map[string]string{
						"uri":       "http://localhost",
						"auth_type": "basic",
						"username":  "some-user-name",
						"password":  "some-password",
						"token":     "some-token",
						"headers":   `bad header`,
					},
				},
			},
			want:    options{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseOptions(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseOptions() got = %v, want %v", got, tt.want)
			}
		})
	}
}
