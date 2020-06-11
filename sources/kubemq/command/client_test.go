package command

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/components/config"
	"github.com/kubemq-hub/components/targets/null"
	"github.com/kubemq-hub/components/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/kubemq-hub/components/targets"
)

func setupClient(ctx context.Context, target targets.Target) (*Client, error) {
	c := New()

	err := c.Init(ctx, config.Metadata{
		Name: "kubemq-rpc",
		Kind: "",
		Properties: map[string]string{
			"host":                       "localhost",
			"port":                       "50000",
			"client_id":                  "",
			"auth_token":                 "some-auth token",
			"channel":                    "command",
			"group":                      "",
			"concurrency":                "1",
			"auto_reconnect":             "true",
			"reconnect_interval_seconds": "1",
			"max_reconnects":             "0",
		},
	})
	if err != nil {
		return nil, err
	}
	err = c.Start(ctx, target)
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second)
	return c, nil
}

func TestClient_processQuery(t *testing.T) {
	tests := []struct {
		name         string
		target       targets.Target
		req          *types.Request
		wantResp     *types.Response
		timeout      time.Duration
		wantQueryErr bool
		wantErr      bool
	}{
		{
			name: "request",
			target: &null.Client{
				Delay:         0,
				DoError:       nil,
				ResponseError: nil,
			},
			req:          types.NewRequest().SetData([]byte("some-data")),
			wantResp:     types.NewResponse().SetData([]byte("some-data")),
			timeout:      5 * time.Second,
			wantQueryErr: false,
			wantErr:      false,
		},
		{
			name: "request with target do error",
			target: &null.Client{
				Delay:         0,
				DoError:       fmt.Errorf("do-error"),
				ResponseError: nil,
			},
			req:          types.NewRequest().SetData([]byte("some-data")),
			wantResp:     types.NewResponse().SetError("do-error"),
			timeout:      5 * time.Second,
			wantQueryErr: true,
			wantErr:      false,
		},
		{
			name: "request with target remote error",
			target: &null.Client{
				Delay:         0,
				DoError:       nil,
				ResponseError: fmt.Errorf("do-error"),
			},
			req:          types.NewRequest().SetData([]byte("some-data")),
			wantResp:     types.NewResponse().SetError("do-error"),
			timeout:      5 * time.Second,
			wantQueryErr: false,
			wantErr:      false,
		},
		{
			name: "bad request",
			target: &null.Client{
				Delay:         0,
				DoError:       nil,
				ResponseError: nil,
			},
			req:          nil,
			wantResp:     nil,
			timeout:      5 * time.Second,
			wantQueryErr: true,
			wantErr:      false,
		},
		{
			name: "request timeout",
			target: &null.Client{
				Delay:         4 * time.Second,
				DoError:       nil,
				ResponseError: nil,
			},
			req:          types.NewRequest().SetData([]byte("some-data")),
			wantResp:     nil,
			timeout:      3 * time.Second,
			wantQueryErr: false,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c, err := setupClient(ctx, tt.target)
			require.NoError(t, err)
			defer func() {
				_ = c.Stop()
			}()
			command := c.client.Q().
				SetChannel("command").
				SetTimeout(tt.timeout).
				SetMetadata("some metadata")
			if tt.req != nil {
				command.SetBody(tt.req.MarshalBinary())
			}
			commandResp, err := command.Send(ctx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, commandResp)
			if tt.wantQueryErr {
				require.NotEmpty(t, commandResp.Error)
				return
			}
			require.NotNil(t, commandResp.Body)
			gotResponse := &types.Response{}
			require.NoError(t, gotResponse.UnmarshalBinary(commandResp.Body))
			require.EqualValues(t, tt.wantResp, gotResponse)
		})
	}
}

func TestClient_Init(t *testing.T) {

	tests := []struct {
		name    string
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Metadata{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"host":                       "localhost",
					"port":                       "50000",
					"client_id":                  "",
					"auth_token":                 "some-auth token",
					"channel":                    "some-channel",
					"group":                      "",
					"concurrency":                "1",
					"auto_reconnect":             "true",
					"reconnect_interval_seconds": "1",
					"max_reconnects":             "0",
				},
			},
			wantErr: false,
		},
		{
			name: "init - error",
			cfg: config.Metadata{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"host": "localhost",
					"port": "-1",
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
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Start(t *testing.T) {

	tests := []struct {
		name    string
		target  targets.Target
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "start",
			target: &null.Client{
				Delay:         0,
				DoError:       nil,
				ResponseError: nil,
			},
			cfg: config.Metadata{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"host":                       "localhost",
					"port":                       "50000",
					"client_id":                  "",
					"auth_token":                 "some-auth token",
					"channel":                    "some-channel",
					"group":                      "",
					"concurrency":                "1",
					"auto_reconnect":             "false",
					"reconnect_interval_seconds": "1",
					"max_reconnects":             "0",
				},
			},
			wantErr: false,
		},
		{
			name:   "start - bad target",
			target: nil,
			cfg: config.Metadata{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"host":                       "localhost",
					"port":                       "50000",
					"client_id":                  "",
					"auth_token":                 "some-auth token",
					"channel":                    "some-channel",
					"group":                      "",
					"concurrency":                "1",
					"auto_reconnect":             "true",
					"reconnect_interval_seconds": "1",
					"max_reconnects":             "0",
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
			_ = c.Init(ctx, tt.cfg)

			if err := c.Start(ctx, tt.target); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
