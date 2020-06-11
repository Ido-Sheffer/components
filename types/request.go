package types

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/kubemq-io/kubemq-go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Request struct {
	Metadata Metadata `json:"metadata"`
	Data     []byte   `json:"data"`
}

func NewRequest() *Request {
	return &Request{}
}

func (r *Request) SetMetadata(value Metadata) *Request {
	r.Metadata = value
	return r
}

func (r *Request) SetData(value []byte) *Request {
	r.Data = value
	return r
}

func ParseRequest(body []byte) (*Request, error) {
	if body == nil {
		return nil, fmt.Errorf("empty request")
	}
	req := &Request{}
	err := json.Unmarshal(body, req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *Request) MarshalBinary() []byte {
	data, _ := json.Marshal(r)
	return data
}

func parseRequest(meta string, body []byte) (*Request, error) {
	req := NewRequest()
	parsedMeta, err := UnmarshallMetadata(meta)
	if err != nil {
		return nil, fmt.Errorf("error parsing request metadata, %w", err)
	}
	return req.
			SetMetadata(parsedMeta).
			SetData(body),
		nil
}
func ParseRequestFromEvent(event *kubemq.Event) (*Request, error) {
	return parseRequest(event.Metadata, event.Body)
}

func ParseRequestFromEventStore(event *kubemq.EventStore) (*Request, error) {
	return parseRequest(event.Metadata, event.Body)
}
func ParseRequestFromEventStoreReceive(event *kubemq.EventStoreReceive) (*Request, error) {
	return parseRequest(event.Metadata, event.Body)
}

func ParseRequestFromCommand(cmd *kubemq.Command) (*Request, error) {
	return parseRequest(cmd.Metadata, cmd.Body)
}
func ParseRequestFromCommandReceive(cmd *kubemq.CommandReceive) (*Request, error) {
	return parseRequest(cmd.Metadata, cmd.Body)
}

func ParseRequestFromQuery(query *kubemq.Query) (*Request, error) {
	return parseRequest(query.Metadata, query.Body)
}

func ParseRequestFromQueryReceive(query *kubemq.QueryReceive) (*Request, error) {
	return parseRequest(query.Metadata, query.Body)
}

func ParseRequestFromQueueMessage(msg *kubemq.QueueMessage) (*Request, error) {
	return parseRequest(msg.Metadata, msg.Body)
}
