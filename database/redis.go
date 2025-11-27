package database

import (
	"github.com/mszlu521/thunder/config"
	"github.com/mszlu521/thunder/db"
)

var (
	RedisCli *db.Redis
)

func InitRedis(redisConf *config.Redis) {
	if redisConf == nil {
		return
	}

	r := db.Redis{}
	r.Init(redisConf)
	RedisCli = &r
}