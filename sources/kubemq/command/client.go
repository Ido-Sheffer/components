package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/kubemq-hub/components/config"
	"github.com/kubemq-hub/components/pkg/logger"
	"github.com/kubemq-hub/components/targets"
	"github.com/kubemq-hub/components/types"
	"github.com/nats-io/nuid"
	"time"

	"github.com/kubemq-io/kubemq-go"
)

const (
	defaultHost          = "localhost"
	defaultPort          = 50000
	defaultMaxReconnects = 0
	defaultAutoReconnect = true
)

var (
	errInvalidTarget = errors.New("invalid target received, cannot be nil")
)

type Client struct {
	name   string
	opts   options
	client *kubemq.Client
	log    *logger.Logger
	target targets.Target
}

func New() *Client {
	return &Client{}

}
func (c *Client) Name() string {
	return c.name
}
func (c *Client) Init(ctx context.Context, cfg config.Metadata) error {
	c.name = cfg.Name
	c.log = logger.NewLogger(cfg.Name)
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	c.client, _ = kubemq.NewClient(ctx,
		kubemq.WithAddress(c.opts.host, c.opts.port),
		kubemq.WithClientId(c.opts.clientId),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithAuthToken(c.opts.authToken),
		kubemq.WithCheckConnection(true),
		kubemq.WithMaxReconnects(c.opts.maxReconnects),
		kubemq.WithAutoReconnect(c.opts.autoReconnect),
		kubemq.WithReconnectInterval(c.opts.reconnectIntervalSeconds))
	return nil
}

func (c *Client) Start(ctx context.Context, target targets.Target) error {
	if target == nil {
		return errInvalidTarget
	} else {
		c.target = target
	}
	group := nuid.Next()
	if c.opts.group != "" {
		group = c.opts.group
	}
	for i := 0; i < c.opts.concurrency; i++ {
		errCh := make(chan error, 1)
		commandsCh, err := c.client.SubscribeToCommands(ctx, c.opts.channel, group, errCh)
		if err != nil {
			return fmt.Errorf("error on subscribing to command channel, %w", err)
		}
		go func(ctx context.Context, commandCh <-chan *kubemq.CommandReceive, errCh chan error) {
			c.run(ctx, commandsCh, errCh)
		}(ctx, commandsCh, errCh)
	}
	return nil
}

func (c *Client) run(ctx context.Context, commandCh <-chan *kubemq.CommandReceive, errCh chan error) {
	for {
		select {
		case command := <-commandCh:
			go func(q *kubemq.CommandReceive) {
				cmdResponse := c.client.R().
					SetRequestId(command.Id).
					SetResponseTo(command.ResponseTo)
				resp, err := c.processCommand(ctx, command)
				if err != nil {
					cmdResponse.SetError(err)
				} else {
					if resp.Error != "" {
						cmdResponse.SetError(fmt.Errorf(resp.Error))
					} else {
						cmdResponse.SetExecutedAt(time.Now())
					}
				}
				err = cmdResponse.Send(ctx)
				if err != nil {
					c.log.Errorf("error sending command response %s", err.Error())
				}
			}(command)

		case err := <-errCh:
			c.log.Errorf("error received from kuebmq server, %s", err.Error())
			return
		case <-ctx.Done():
			return

		}
	}
}

func (c *Client) processCommand(ctx context.Context, command *kubemq.CommandReceive) (*types.Response, error) {
	req, err := types.ParseRequestFromCommandReceive(command)
	if err != nil {
		return nil, fmt.Errorf("invalid request format, %w", err)
	}
	resp, err := c.target.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *Client) Stop() error {
	return c.client.Close()
}
