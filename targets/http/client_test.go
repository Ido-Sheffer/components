package http

import (
	"context"
	"github.com/kubemq-hub/components/config"
	"github.com/kubemq-hub/components/types"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func startHttpTestServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, values := range r.Header {
			w.Header().Set(key, strings.Join(values, ","))
		}
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(err)
			return
		}
		_, err = w.Write(buf)
		if err != nil {
			log.Error(err)
		}

	}))
}

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Metadata
		request *types.Request
		want    *types.Response
		wantErr bool
	}{
		{
			name: "valid request",
			cfg: config.Metadata{
				Name: "target.http",
				Kind: "target.http",
				Properties: map[string]string{
					"uri":       "",
					"auth_type": "no_auth",
					"username":  "",
					"password":  "",
					"token":     "",
					"headers":   ``,
				},
			},
			request: types.NewRequest().
				SetMethod("post").
				SetData([]byte("some-data")),
			want: types.NewResponse().
				SetData([]byte("some-data")),

			wantErr: false,
		},
		{
			name: "invalid chore - bad method",
			cfg: config.Metadata{
				Name: "target.http",
				Kind: "target.http",
				Properties: map[string]string{
					"uri":       "",
					"auth_type": "no_auth",
					"username":  "",
					"password":  "",
					"token":     "",
					"headers":   ``,
				},
			},
			request: types.NewRequest().
				SetMethod("invalid").
				SetData([]byte("some-data")),
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid chore - bad url",
			cfg: config.Metadata{
				Name: "target.http",
				Kind: "target.http",
				Properties: map[string]string{
					"uri":       "bad-url",
					"auth_type": "no_auth",
					"username":  "",
					"password":  "",
					"token":     "",
					"headers":   ``,
				},
			},
			request: types.NewRequest().
				SetMethod("post").
				SetData([]byte("some-data")),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			require.NoError(t, err)
			ts := startHttpTestServer(t)
			defer ts.Close()
			tt.request.Url = ts.URL
			got, err := c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			require.EqualValues(t, got.Data, tt.want.Data)
		})
	}
}
