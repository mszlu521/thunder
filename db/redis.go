package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"thunder/config"
	"time"
)

type Redis struct {
	Options *redis.Options
	Client  *redis.Client
}

func (r *Redis) Init() {
	if r.Options == nil {
		r.Options = &redis.Options{
			Addr:           config.Conf.DB.Redis.Addr,
			DB:             config.Conf.DB.Redis.DB,
			Password:       config.Conf.DB.Redis.Password,
			PoolSize:       config.Conf.DB.Redis.PoolSize,
			MaxIdleConns:   config.Conf.DB.Redis.MaxIdleConns,
			MaxActiveConns: config.Conf.DB.Redis.MaxOpenConns,
		}
	}
	rdb := redis.NewClient(r.Options)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	r.Client = rdb
}
