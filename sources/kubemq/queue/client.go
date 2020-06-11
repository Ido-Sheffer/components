package queue

import (
	"context"
	"github.com/kubemq-hub/components/targets"
	"github.com/kubemq-hub/components/types"
	"github.com/kubemq-io/kubemq-go"

	"errors"
	"fmt"
	"github.com/kubemq-hub/components/config"
	"github.com/kubemq-hub/components/pkg/logger"
)

var (
	errInvalidTarget = errors.New("invalid controller received, cannot be null")
)

const (
	defaultHost        = "localhost"
	defaultPort        = 50000
	defaultBatchSize   = 1
	defaultWaitTimeout = 60
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
	c.log = logger.NewLogger(fmt.Sprintf("kubemq-evetns-source-%s", cfg.Name))
	var err error
	c.opts, err = parseOptions(cfg)
	if err != nil {
		return err
	}
	c.client, err = kubemq.NewClient(ctx,
		kubemq.WithAddress(c.opts.host, c.opts.port),
		kubemq.WithClientId(c.opts.clientId),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithAuthToken(c.opts.authToken))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Start(ctx context.Context, target targets.Target) error {
	if target == nil {
		return errInvalidTarget
	} else {
		c.target = target
	}
	for i := 0; i < c.opts.concurrency; i++ {
		go c.run(ctx)
	}
	return nil
}

func (c *Client) run(ctx context.Context) {
	for {
		queueMessages, err := c.getQueueMessages(ctx)
		if err != nil {
			c.log.Error(err.Error())
			return
		}
		for _, message := range queueMessages {
			err := c.processQueueMessage(ctx, message)
			if err != nil {
				c.log.Error(err.Error())
			}
		}
		select {
		case <-ctx.Done():
			return
		default:

		}
	}

}
func (c *Client) getQueueMessages(ctx context.Context) ([]*kubemq.QueueMessage, error) {
	receiveResult, err := c.client.NewReceiveQueueMessagesRequest().
		SetChannel(c.opts.channel).
		SetMaxNumberOfMessages(c.opts.batchSize).
		SetWaitTimeSeconds(c.opts.waitTimeout).
		Send(ctx)
	if err != nil {
		return nil, err
	}
	return receiveResult.Messages, nil
}

func (c *Client) processQueueMessage(ctx context.Context, msg *kubemq.QueueMessage) error {
	response := newResponse(c.client)
	req, err := types.ParseRequest(msg.Body)
	if err != nil {
		return response.SetError(fmt.Errorf("invalid request format, %w", err)).Send(ctx)
	}
	if req.ResponseQueue != "" {
		response.SetChannel(req.ResponseQueue)
	} else {
		response.SetChannel(c.opts.responseChannel)
	}
	r, err := c.target.Do(ctx, req)
	if err != nil {
		response.SetBody(types.NewResponse().SetError(err.Error()).MarshalBinary())
	} else {
		response.SetBody(r.MarshalBinary())
	}
	return response.Send(ctx)

}

func (c *Client) Stop() error {
	return c.client.Close()
}
