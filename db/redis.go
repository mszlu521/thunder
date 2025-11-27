package db

import (
	"context"
	"github.com/mszlu521/thunder/config"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	Options *redis.Options
	Client  *redis.Client
}

func (r *Redis) Init(redisConf *config.Redis) {
	if r.Options == nil {
		r.Options = &redis.Options{
			Addr:         redisConf.GetAddr(),
			DB:           redisConf.GetDB(),
			Password:     redisConf.GetPassword(),
			PoolSize:     redisConf.GetPoolSize(),
			MaxIdleConns: redisConf.GetMaxIdleConns(),
			MaxActiveConns: redisConf.GetMaxOpenConns(),
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