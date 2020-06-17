package mongodb

import (
	"fmt"
	"github.com/kubemq-hub/components/config"
	"math"
	"time"
)

var writeConcurrencyMap = map[string]string{
	"majority": "majority",
	"":         "",
}

var readConcurrencyMap = map[string]string{
	"local":        "local",
	"majority":     "majority",
	"available":    "available",
	"linearizable": "linearizable",
	"snapshot":     "snapshot",
	"":             "",
}

type options struct {
	host             string
	username         string
	password         string
	database         string
	collection       string
	writeConcurrency string
	readConcurrency  string
	params           string
	operationTimeout time.Duration
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{}
	var err error
	o.host, err = cfg.MustParseString("host")
	if err != nil {
		return options{}, fmt.Errorf("error parsing host, %w", err)
	}
	o.username = cfg.ParseString("username", "")
	o.password = cfg.ParseString("password", "")
	o.database, err = cfg.MustParseString("database")
	if err != nil {
		return options{}, fmt.Errorf("error parsing database name, %w", err)
	}
	o.collection, err = cfg.MustParseString("collection")
	if err != nil {
		return options{}, fmt.Errorf("error parsing collection name, %w", err)
	}

	o.writeConcurrency, err = cfg.ParseStringMap("write_concurrency", writeConcurrencyMap)
	if err != nil {
		return options{}, fmt.Errorf("error on parsing write concurrency, %w", err)
	}

	o.readConcurrency, err = cfg.ParseStringMap("read_concurrency", readConcurrencyMap)
	if err != nil {
		return options{}, fmt.Errorf("error on parsing read concurrency, %w", err)
	}

	o.params = cfg.ParseString("params", "")
	operationTimeoutSeconds, err := cfg.ParseIntWithRange("operation_timeout_seconds", 2, 0, math.MaxInt32)
	if err != nil {
		return options{}, fmt.Errorf("error operation timeout seconds, %w", err)
	}
	o.operationTimeout = time.Duration(operationTimeoutSeconds) * time.Second
	return o, nil
}
