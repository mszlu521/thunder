package cache

import (
	"context"
	"errors"
	"github.com/mszlu521/thunder/database"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	redisCli *redis.Client
}

func (c *RedisCache) Get(key string) (string, error) {
	result, err := c.redisCli.Get(context.Background(), key).Result()
	if err != nil && errors.Is(err, redis.Nil) {
		return "", nil
	}
	return result, err
}
func (c *RedisCache) TTL(key string) (time.Duration, error) {
	result, err := c.redisCli.TTL(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
	}
	return result, err
}
func (c *RedisCache) GetValueAndTTL(ctx context.Context, key string) (string, time.Duration, error) {
	rdb := c.redisCli
	pipe := rdb.Pipeline()
	getCmd := pipe.Get(ctx, key)
	ttlCmd := pipe.TTL(ctx, key)
	_, err := pipe.Exec(ctx)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", 0, nil
		}
		return "", 0, err
	}
	val, err := getCmd.Result()
	if err != nil {
		return "", 0, err
	}

	ttl, err := ttlCmd.Result()
	if err != nil {
		return "", 0, err
	}
	return val, ttl, nil
}

func (c *RedisCache) Set(key string, value string, expire int64) error {
	return c.redisCli.Set(context.Background(), key, value, time.Duration(expire)*time.Second).Err()
}

func (c *RedisCache) Exist(key string) bool {
	result, err := c.redisCli.Exists(context.Background(), key).Result()
	return result == 1 && err == nil
}
func NewRedisCache() *RedisCache {
	return &RedisCache{
		redisCli: database.RedisCli.Client,
	}
}
