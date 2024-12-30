package database

import (
	"github.com/redis/go-redis/v9"
	"thunder/config"
	"thunder/db"
)

var (
	RedisCli *db.Redis
)

func InitRedis() {
	r := db.Redis{
		Options: &redis.Options{
			Addr:           config.Conf.DB.Redis.Addr,
			DB:             config.Conf.DB.Redis.DB,
			Password:       config.Conf.DB.Redis.Password,
			MaxActiveConns: config.Conf.DB.Redis.MaxOpenConns,
			PoolSize:       config.Conf.DB.Redis.PoolSize,
			MaxIdleConns:   config.Conf.DB.Redis.MaxIdleConns,
		},
	}
	r.Init()
	RedisCli = &r
}
