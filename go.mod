module github.com/kubemq-hub/components

go 1.14

require (
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/coreos/etcd v3.3.22+incompatible // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-redis/redis/v7 v7.4.0
	github.com/go-resty/resty/v2 v2.3.0
	github.com/json-iterator/go v1.1.10
	github.com/kubemq-io/kubemq-go v1.4.0
	github.com/labstack/echo/v4 v4.1.16
	github.com/nats-io/nuid v1.0.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.4.0
	go.etcd.io/etcd v3.3.22+incompatible
	go.mongodb.org/mongo-driver v1.3.4
	go.uber.org/zap v1.10.0
	google.golang.org/grpc v1.28.0
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace github.com/kubemq-io/kubemq-go => ../../kubemq/kubemq-go
replace (
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
