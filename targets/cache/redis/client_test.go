package redis

import (
	"context"
	"github.com/kubemq-hub/components/config"
	"github.com/kubemq-hub/components/types"
	"github.com/nats-io/nuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClient_Init(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Metadata{
				Name: "redis-target",
				Kind: "",
				Properties: map[string]string{
					"host":                        "localhost:6379",
					"password":                    "",
					"enable_tls":                  "false",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
				},
			},
			wantErr: false,
		},
		{
			name: "init - error no host",
			cfg: config.Metadata{
				Name: "redis-target",
				Kind: "",
				Properties: map[string]string{

					"password":                    "",
					"enable_tls":                  "false",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
				},
			},
			wantErr: true,
		},
		{
			name: "init - error",
			cfg: config.Metadata{
				Name: "redis-target",
				Kind: "",
				Properties: map[string]string{
					"host":                        "localhost:2000",
					"password":                    "",
					"enable_tls":                  "false",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
				},
			},
			wantErr: true,
		}, {
			name: "init - bad options - 1",
			cfg: config.Metadata{
				Name: "redis-target",
				Kind: "",
				Properties: map[string]string{
					"host":                        "localhost:6397",
					"password":                    "",
					"enable_tls":                  "false",
					"max_retries":                 "-1",
					"max_retries_backoff_seconds": "0",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad options - 2",
			cfg: config.Metadata{
				Name: "redis-target",
				Kind: "",
				Properties: map[string]string{
					"host":                        "localhost:6397",
					"password":                    "",
					"enable_tls":                  "false",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "-1",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()

			if err := c.Init(ctx, tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantSetErr %v", err, tt.wantErr)
				return
			}
			require.EqualValues(t, tt.cfg.Name, c.Name())
		})
	}
}
func TestClient_Set_Get(t *testing.T) {
	tests := []struct {
		name            string
		cfg             config.Metadata
		setRequest      *types.Request
		getRequest      *types.Request
		wantSetResponse *types.Response
		wantGetResponse *types.Response
		wantSetErr      bool
		wantGetErr      bool
	}{
		{
			name: "valid set get request",
			cfg: config.Metadata{
				Name: "target.redis",
				Kind: "target.redis",
				Properties: map[string]string{
					"host":                        "localhost:6379",
					"password":                    "",
					"enable_tls":                  "false",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-key").
				SetData([]byte("some-data")),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("key", "some-key"),

			wantSetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("result", "ok"),
			wantGetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("error", "false").
				SetData([]byte("some-data")),
			wantSetErr: false,
			wantGetErr: false,
		},
		{
			name: "valid set , no key get request",
			cfg: config.Metadata{
				Name: "target.redis",
				Kind: "target.redis",
				Properties: map[string]string{
					"host":                        "localhost:6379",
					"password":                    "",
					"enable_tls":                  "false",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
				},
			},
			setRequest: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-key").
				SetData([]byte("some-data")),
			getRequest: types.NewRequest().
				SetMetadataKeyValue("method", "get").
				SetMetadataKeyValue("key", "bad-key"),

			wantSetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("result", "ok"),
			wantGetResponse: types.NewResponse().
				SetMetadataKeyValue("key", "bad-key").
				SetMetadataKeyValue("error", "true").
				SetMetadataKeyValue("message", "no data found for this key"),
			wantSetErr: false,
			wantGetErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			c := New()
			err := c.Init(ctx, tt.cfg)
			require.NoError(t, err)
			gotSetResponse, err := c.Do(ctx, tt.setRequest)
			if tt.wantSetErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotSetResponse)
			require.EqualValues(t, tt.wantSetResponse, gotSetResponse)
			gotGetResponse, err := c.Do(ctx, tt.getRequest)
			if tt.wantGetErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, gotGetResponse)
			require.EqualValues(t, tt.wantGetResponse, gotGetResponse)
		})
	}
}
func TestClient_Delete(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := New()
	err := c.Init(ctx, config.Metadata{
		Name: "target.redis",
		Kind: "target.redis",
		Properties: map[string]string{
			"host":                        "localhost:6379",
			"password":                    "",
			"enable_tls":                  "false",
			"max_retries":                 "0",
			"max_retries_backoff_seconds": "0",
		},
	})
	key := nuid.Next()
	require.NoError(t, err)
	setRequest := types.NewRequest().
		SetMetadataKeyValue("method", "set").
		SetMetadataKeyValue("key", key).
		SetData([]byte("some-data"))

	_, err = c.Do(ctx, setRequest)
	require.NoError(t, err)
	getRequest := types.NewRequest().
		SetMetadataKeyValue("method", "get").
		SetMetadataKeyValue("key", key)
	gotGetResponse, err := c.Do(ctx, getRequest)
	require.NoError(t, err)
	require.NotNil(t, gotGetResponse)
	require.EqualValues(t, []byte("some-data"), gotGetResponse.Data)

	delRequest := types.NewRequest().
		SetMetadataKeyValue("method", "delete").
		SetMetadataKeyValue("key", key)
	_, err = c.Do(ctx, delRequest)
	require.NoError(t, err)
	gotGetResponse, err = c.Do(ctx, getRequest)
	require.NoError(t, err)
	require.NotNil(t, gotGetResponse)
	require.EqualValues(t, []byte(nil), gotGetResponse.Data)
}
func TestClient_Do(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Metadata
		request *types.Request
		wantErr bool
	}{
		{
			name: "valid request",
			cfg: config.Metadata{
				Name: "target.redis",
				Kind: "target.redis",
				Properties: map[string]string{
					"host":                        "localhost:6379",
					"password":                    "",
					"enable_tls":                  "false",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-key").
				SetData([]byte("some-data")),
			wantErr: false,
		},
		{
			name: "invalid request - bad method",
			cfg: config.Metadata{
				Name: "target.redis",
				Kind: "target.redis",
				Properties: map[string]string{
					"host":                        "localhost:6379",
					"password":                    "",
					"enable_tls":                  "false",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "bad-method").
				SetMetadataKeyValue("key", "some-key").
				SetData([]byte("some-data")),
			wantErr: true,
		},
		{
			name: "invalid request - no key",
			cfg: config.Metadata{
				Name: "target.redis",
				Kind: "target.redis",
				Properties: map[string]string{
					"host":                        "localhost:6379",
					"password":                    "",
					"enable_tls":                  "false",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetData([]byte("some-data")),
			wantErr: true,
		},
		{
			name: "invalid request - bad etag",
			cfg: config.Metadata{
				Name: "target.redis",
				Kind: "target.redis",
				Properties: map[string]string{
					"host":                        "localhost:6379",
					"password":                    "",
					"enable_tls":                  "false",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("etag", "-1").
				SetData([]byte("some-data")),
			wantErr: true,
		},
		{
			name: "invalid request - bad concurrency",
			cfg: config.Metadata{
				Name: "target.redis",
				Kind: "target.redis",
				Properties: map[string]string{
					"host":                        "localhost:6379",
					"password":                    "",
					"enable_tls":                  "false",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("concurrency", "bad-concurrency").
				SetData([]byte("some-data")),
			wantErr: true,
		}, {
			name: "invalid request - bad consistency",
			cfg: config.Metadata{
				Name: "target.redis",
				Kind: "target.redis",
				Properties: map[string]string{
					"host":                        "localhost:6379",
					"password":                    "",
					"enable_tls":                  "false",
					"max_retries":                 "0",
					"max_retries_backoff_seconds": "0",
				},
			},
			request: types.NewRequest().
				SetMetadataKeyValue("method", "set").
				SetMetadataKeyValue("key", "some-key").
				SetMetadataKeyValue("consistency", "bad-consistency").
				SetData([]byte("some-data")),
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
			_, err = c.Do(ctx, tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}
}
