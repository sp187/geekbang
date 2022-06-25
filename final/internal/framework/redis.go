package fw

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type AppCache interface {
	Get(context.Context, string) (string, error)
	Delete(context.Context, string) error
	Set(context.Context, string, interface{}, time.Duration) error
}

type Redis struct {
	client *redis.Client
}

func NewRedisClient(address, password string) AppCache {
	return &Redis{
		client: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
		}),
	}
}

func (cache *Redis) Get(ctx context.Context, key string) (string, error) {
	return cache.client.Get(ctx, key).Result()
}

func (cache *Redis) Delete(ctx context.Context, key string) error {
	return cache.client.Del(ctx, key).Err()
}

func (cache *Redis) Set(ctx context.Context, key string, data interface{}, expiration time.Duration) error {
	infoByte, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = cache.client.Set(ctx, key, string(infoByte), expiration).Err()
	return err
}
