package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go-micro.dev/v4/cache"
	"go-micro.dev/v4/util/cmd"
)

func init() {
	cmd.DefaultCaches["redis"] = NewCache
}

// NewCache returns a new redis cache.
func NewCache(opts ...cache.Option) cache.Cache {
	options := cache.NewOptions(opts...)

	if options.Context == nil {
		options.Context = context.Background()
	}

	return &redisCache{
		opts:   options,
		client: options.Context.Value(redisClientKey{}).(redis.Cmdable),
	}
}

type redisCache struct {
	opts   cache.Options
	client redis.Cmdable
}

func (c *redisCache) Get(ctx context.Context, key string) (interface{}, time.Time, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if err != nil && err == redis.Nil {
		return nil, time.Time{}, cache.ErrKeyNotFound
	} else if err != nil {
		return nil, time.Time{}, err
	}

	dur, err := c.client.TTL(ctx, key).Result()
	if err != nil {
		return nil, time.Time{}, err
	}
	if dur == -1 {
		return val, time.Unix(1<<63-1, 0), nil
	}
	if dur == -2 {
		return val, time.Time{}, cache.ErrItemExpired
	}

	return val, time.Now().Add(dur), nil
}

func (c *redisCache) Put(ctx context.Context, key string, val interface{}, dur time.Duration) error {
	return c.client.Set(ctx, key, val, dur).Err()
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (m *redisCache) String() string {
	return "redis"
}
