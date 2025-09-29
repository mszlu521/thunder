package database

import (
	"github.com/mszlu521/thunder/config"
	"github.com/mszlu521/thunder/db"
)

var (
	_pg *db.Postgres
)

func InitPostgres(pgConf config.Postgres) {
	p := db.Postgres{
		Database:     pgConf.Database,
		Host:         pgConf.Host,
		MaxIdleConns: pgConf.MaxIdleConns,
		MaxOpenConns: pgConf.MaxOpenConns,
		Password:     pgConf.Password,
		Port:         pgConf.Port,
		Username:     pgConf.User,
		SSLMode:      pgConf.SSLMode,
		PingTimeout:  pgConf.PingTimeout,
	}
	err := p.Init()
	if err != nil {
		panic(err)
	}
	_pg = &p
}

func GetPostgresDB() *db.Postgres {
	return _pg
}