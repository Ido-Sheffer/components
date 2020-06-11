package queue

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/components/config"
	"github.com/kubemq-hub/components/targets/null"
	"github.com/kubemq-hub/components/types"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/kubemq-hub/components/targets"
)

func setupClient(ctx context.Context, target targets.Target) (*Client, error) {
	c := New()

	err := c.Init(ctx, config.Metadata{
		Name: "kubemq-queue",
		Kind: "",
		Properties: map[string]string{
			"host":             "localhost",
			"port":             "50000",
			"client_id":        "some-client-id",
			"auth_token":       "",
			"channel":          "queue",
			"response_channel": "default-response",
			"concurrency":      "1",
			"batch_size":       "1",
			"wait_timeout":     "60",
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

func getResponseFromQueue(ctx context.Context, c *Client, channel string, timeout int) (*kubemq.QueueMessage, error) {
	respMsgs, err := c.client.ReceiveQueueMessages(ctx,
		c.client.NewReceiveQueueMessagesRequest().
			SetChannel(channel).
			SetClientId(nuid.Next()).
			SetMaxNumberOfMessages(1).
			SetWaitTimeSeconds(timeout))
	if err != nil {
		return nil, err
	}
	if len(respMsgs.Messages) == 0 {
		return nil, fmt.Errorf("no messages")
	}
	return respMsgs.Messages[0], nil

}
func TestClient_processQueue(t *testing.T) {
	tests := []struct {
		name            string
		target          targets.Target
		respChannel     string
		req             *types.Request
		wantResp        *types.Response
		timeout         int
		wantResponseErr bool
		wantErr         bool
	}{
		{
			name: "request",
			target: &null.Client{
				Delay:         0,
				DoError:       nil,
				ResponseError: nil,
			},
			respChannel:     "queue.response.1",
			req:             types.NewRequest().SetData([]byte("some-data")).SetResponseQueue("queue.response.1"),
			wantResp:        types.NewResponse().SetData([]byte("some-data")),
			timeout:         5,
			wantResponseErr: false,
			wantErr:         false,
		},
		{
			name: "request with default response channel",
			target: &null.Client{
				Delay:         0,
				DoError:       nil,
				ResponseError: nil,
			},
			respChannel:     "default-response",
			req:             types.NewRequest().SetData([]byte("some-data")),
			wantResp:        types.NewResponse().SetData([]byte("some-data")),
			timeout:         5,
			wantResponseErr: false,
			wantErr:         false,
		},
		{
			name: "request with target do error",
			target: &null.Client{
				Delay:         0,
				DoError:       fmt.Errorf("do-error"),
				ResponseError: nil,
			},
			respChannel:     "queue.response.2",
			req:             types.NewRequest().SetData([]byte("some-data")).SetResponseQueue("queue.response.2"),
			wantResp:        types.NewResponse().SetError("do-error"),
			timeout:         5,
			wantResponseErr: false,
			wantErr:         false,
		},
		{
			name: "request with target remote error",
			target: &null.Client{
				Delay:         0,
				DoError:       nil,
				ResponseError: fmt.Errorf("do-error"),
			},
			respChannel:     "queue.response.3",
			req:             types.NewRequest().SetData([]byte("some-data")).SetResponseQueue("queue.response.3"),
			wantResp:        types.NewResponse().SetError("do-error"),
			timeout:         5,
			wantResponseErr: false,
			wantErr:         false,
		},
		{
			name: "bad request",
			target: &null.Client{
				Delay:         0,
				DoError:       nil,
				ResponseError: nil,
			},
			req:             nil,
			wantResp:        nil,
			timeout:         5,
			wantResponseErr: true,
			wantErr:         false,
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
			msg := c.client.NewQueueMessage().
				SetChannel("queue").
				SetMetadata("some metadata")
			if tt.req != nil {
				msg.SetBody(tt.req.MarshalBinary())
			}
			msgResp, err := msg.Send(ctx)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, msgResp)

			resQueueMessage, err := getResponseFromQueue(ctx, c, tt.respChannel, tt.timeout)
			if tt.wantResponseErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resQueueMessage.Body)
			gotResponse := &types.Response{}
			require.NoError(t, gotResponse.UnmarshalBinary(resQueueMessage.Body))
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
					"host":             "localhost",
					"port":             "50000",
					"client_id":        "some-client-id",
					"auth_token":       "some-auth token",
					"channel":          "some-channel",
					"response_channel": "some-response-channel",
					"concurrency":      "1",
					"batch_size":       "1",
					"wait_timeout":     "60",
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
					"host":             "localhost",
					"port":             "50000",
					"client_id":        "some-client-id",
					"auth_token":       "some-auth token",
					"channel":          "some-channel",
					"response_channel": "some-response-channel",
					"concurrency":      "1",
					"batch_size":       "1",
					"wait_timeout":     "60",
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
					"host":             "localhost",
					"port":             "50000",
					"client_id":        "some-client-id",
					"auth_token":       "some-auth token",
					"channel":          "some-channel",
					"response_channel": "some-response-channel",
					"concurrency":      "1",
					"batch_size":       "1",
					"wait_timeout":     "60",
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