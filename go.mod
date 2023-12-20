module github.com/lenny-mo/payment

go 1.13

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/golang/protobuf v1.5.3
	github.com/google/uuid v1.4.0
	github.com/lenny-mo/emall-utils v0.0.0-20231218141407-3b3960e96cd9
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/config/source/consul/v2 v2.9.1
	github.com/micro/go-plugins/registry/consul/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/trace/opentracing/v2 v2.9.1
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/nats-io/nats.go v1.31.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.5.1
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0
	google.golang.org/protobuf v1.31.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gorm.io/driver/mysql v1.5.2
	gorm.io/gorm v1.25.5
)
