package memcached

import (
	"context"
	"errors"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"

	"github.com/kubemq-hub/components/config"
	"github.com/kubemq-hub/components/pkg/logger"
	"github.com/kubemq-hub/components/types"

	"time"
)

// Client is a Client state store
type Client struct {
	name   string
	client *memcache.Client
	opts   options
	log    *logger.Logger
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

	c.client = memcache.New(c.opts.hosts...)
	c.client.Timeout = time.Duration(c.opts.defaultTimeoutSeconds) * time.Second
	c.client.MaxIdleConns = c.opts.maxIdleConnections
	err = c.client.Ping()
	if err != nil {
		return err
	}
	return nil
}
func (c *Client) Do(ctx context.Context, req *types.Request) (*types.Response, error) {
	meta, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}
	switch meta.method {
	case "get":
		return c.Get(ctx, meta)
	case "set":
		return c.Set(ctx, meta, req.Data)
	case "delete":
		return c.Delete(ctx, meta)

	}
	return nil, nil
}

func (c *Client) Get(ctx context.Context, meta metadata) (*types.Response, error) {
	item, err := c.client.Get(meta.key)
	if err != nil {
		// Return nil for status 204
		if errors.Is(err, memcache.ErrCacheMiss) {
			return types.NewResponse().
				SetMetadataKeyValue("key", meta.key).
				SetMetadataKeyValue("error", "true").
				SetMetadataKeyValue("message", "no data found for this key"), nil
		}
		return nil, err
	}
	return types.NewResponse().
		SetData(item.Value).
		SetMetadataKeyValue("error", "false").
		SetMetadataKeyValue("key", meta.key), nil

}

func (c *Client) Set(ctx context.Context, meta metadata, value []byte) (*types.Response, error) {
	err := c.client.Set(&memcache.Item{Key: meta.key, Value: value})
	if err != nil {
		return nil, fmt.Errorf("failed to set key %s: %s", meta.key, err)
	}
	return types.NewResponse().
			SetMetadataKeyValue("key", meta.key).
			SetMetadataKeyValue("result", "ok"),
		nil
}

func (c *Client) Delete(ctx context.Context, meta metadata) (*types.Response, error) {
	err := c.client.Delete(meta.key)
	if err != nil {
		return nil, fmt.Errorf("failed to delete key '%s',%w", meta.key, err)
	}
	return types.NewResponse().
			SetMetadataKeyValue("key", meta.key).
			SetMetadataKeyValue("result", "ok"),
		nil
}
