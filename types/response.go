package types

import "github.com/kubemq-io/kubemq-go"

type Response struct {
	Metadata Metadata `json:"metadata"`
	Data     []byte   `json:"data"`
	Error    string   `json:"error"`
}

func NewResponse() *Response {
	return &Response{
		Metadata: NewMetadata(),
		Data:     nil,
		Error:    "",
	}
}

func (r *Response) SetMetadata(value Metadata) *Response {
	r.Metadata = value
	return r
}

func (r *Response) SetError(value string) *Response {
	r.Error = value
	return r
}

func (r *Response) SetData(value []byte) *Response {
	r.Data = value
	return r
}

func (r *Response) ToEvent() *kubemq.Event {
	return kubemq.NewEvent().
		SetMetadata(r.Metadata.String()).
		SetBody(r.Data)
}
func (r *Response) ToEventStore() *kubemq.EventStore {
	return kubemq.NewEventStore().
		SetMetadata(r.Metadata.String()).
		SetBody(r.Data)
}

func (r *Response) ToCommand() *kubemq.Command {
	return kubemq.NewCommand().
		SetMetadata(r.Metadata.String()).
		SetBody(r.Data)
}

func (r *Response) ToQuery() *kubemq.Query {
	return kubemq.NewQuery().
		SetMetadata(r.Metadata.String()).
		SetBody(r.Data)
}

func (r *Response) ToQueueMessage() *kubemq.QueueMessage {
	return kubemq.NewQueueMessage().
		SetMetadata(r.Metadata.String()).
		SetBody(r.Data)
}
func (r *Response) ToResponse() *kubemq.Response {
	return kubemq.NewResponse().
		SetMetadata(r.Metadata.String()).
		SetBody(r.Data)
}
