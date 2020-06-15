package command

import (
	"fmt"
	"github.com/kubemq-hub/components/types"
	"math"
	"time"
)

type metadata struct {
	id      string
	channel string
	timeout time.Duration
}

func parseMetadata(meta types.Metadata, opts options) (metadata, error) {
	m := metadata{}
	m.id = meta.ParseString("id", "")
	m.channel = meta.ParseString("channel", opts.defaultChannel)
	timout, err := meta.ParseIntWithRange("timeout_seconds", opts.defaultTimeoutSeconds, 1, math.MaxInt32)
	if err != nil {
		return metadata{}, fmt.Errorf("error ")
	}
	m.timeout = time.Duration(timout) * time.Second
	return m, nil
}
