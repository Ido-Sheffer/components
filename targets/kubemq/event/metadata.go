package event

import (
	"fmt"
	"github.com/kubemq-hub/components/types"
)

type metadata struct {
	id      string
	channel string
}

func parseMetadata(meta types.Metadata, opts options) (metadata, error) {
	m := metadata{}
	m.id = meta.ParseString("id", "")
	m.channel = meta.ParseString("channel", opts.defaultChannel)
	if m.channel == "" {
		return metadata{}, fmt.Errorf("channel cannot be empty")
	}
	return m, nil
}
