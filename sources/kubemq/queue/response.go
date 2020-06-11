package queue

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-go"
)

type Response struct {
	c       *kubemq.Client
	err     error
	body    []byte
	channel string
}

func newResponse(c *kubemq.Client) *Response {
	return &Response{
		c: c,
	}
}

func (r *Response) Send(ctx context.Context) error {
	if r.channel == "" {
		return nil
	}
	qm := r.c.NewQueueMessage().SetChannel(r.channel)
	sendResult, err := qm.SetBody(r.body).Send(ctx)
	if err != nil {
		return err
	}
	if sendResult.IsError {
		return fmt.Errorf(sendResult.Error)
	}
	return nil
}

func (r *Response) SetError(value error) *Response {
	r.err = value
	return r
}
func (r *Response) SetBody(value []byte) *Response {
	r.body = value
	return r
}

func (r *Response) SetChannel(value string) *Response {
	r.channel = value
	return r
}
