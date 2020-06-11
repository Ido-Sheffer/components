package http

import (
	"fmt"
	"github.com/kubemq-hub/components/config"
)

type options struct {
	uri      string
	authType string
	username string
	password string
	token    string
	headers  map[string]string
}

func parseOptions(cfg config.Metadata) (options, error) {
	o := options{
		uri:      "",
		authType: "",
		username: "",
		password: "",
		token:    "",
		headers:  map[string]string{},
	}
	var err error
	o.uri = cfg.ParseString("uri", "")
	o.authType = cfg.ParseString("auth_type", "no_auth")
	o.username = cfg.ParseString("username", "")
	o.password = cfg.ParseString("password", "")
	o.token = cfg.ParseString("token", "")
	o.headers, err = cfg.MustParseJsonMap("headers")
	if err != nil {
		return options{}, fmt.Errorf("error parsing headers value, %w", err)
	}
	return o, nil
}
