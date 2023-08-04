package etcd

import (
	"context"

	"go-micro.dev/v4/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type clientKey struct{}

type logConfigKey struct{}

func Client(client *clientv3.Client) registry.Option {
	return func(o *registry.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, clientKey{}, client)
	}
}

// LogConfig allows you to set etcd log config.
func LogConfig(config *zap.Config) registry.Option {
	return func(o *registry.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, logConfigKey{}, config)
	}
}
