module github.com/kubemq-hub/components

go 1.14

require (
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/ghodss/yaml v1.0.0
	github.com/go-redis/redis/v7 v7.4.0
	github.com/go-resty/resty/v2 v2.3.0
	github.com/json-iterator/go v1.1.10
	github.com/kubemq-io/kubemq-go v1.4.0
	github.com/labstack/echo/v4 v4.1.16
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.4.0
	go.mongodb.org/mongo-driver v1.3.4 // indirect
	go.uber.org/zap v1.10.0
)

replace github.com/kubemq-io/kubemq-go => ../../kubemq/kubemq-go
