module github.com/kubemq-hub/components

go 1.14

require (
	github.com/ghodss/yaml v1.0.0
	github.com/go-resty/resty/v2 v2.0.0
	github.com/json-iterator/go v1.1.10
	github.com/kubemq-io/kubemq-go v1.3.8
	github.com/nats-io/nuid v1.0.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	go.uber.org/zap v1.10.0
)

replace github.com/kubemq-io/kubemq-go => ../../kubemq/kubemq-go
