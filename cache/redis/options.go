package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go-micro.dev/v4/cache"
)

type redisClientKey struct{}

// WithRedisOptions sets advanced options for redis.
func WithRedisClient(client redis.Cmdable) cache.Option {
	return func(o *cache.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, redisClientKey{}, client)
	}
}
