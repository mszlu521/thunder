package database

import (
	"github.com/mszlu521/thunder/config"
	"github.com/mszlu521/thunder/db"
	"github.com/redis/go-redis/v9"
)

var (
	RedisCli *db.Redis
)

func InitRedis(redisConf config.Redis) {
	r := db.Redis{
		Options: &redis.Options{
			Addr:           redisConf.Addr,
			DB:             redisConf.DB,
			Password:       redisConf.Password,
			MaxActiveConns: redisConf.MaxOpenConns,
			PoolSize:       redisConf.PoolSize,
			MaxIdleConns:   redisConf.MaxIdleConns,
		},
	}
	r.Init(redisConf)
	RedisCli = &r
}
