package database

import (
	"github.com/mszlu521/thunder/config"
	"github.com/mszlu521/thunder/db"
)

var (
	_pg *db.Postgres
)

func InitPostgres(pgConf *config.Postgres) {
	if pgConf == nil {
		return
	}

	p := db.Postgres{
		Database:                  pgConf.GetDatabase(),
		Host:                      pgConf.GetHost(),
		MaxIdleConns:              pgConf.GetMaxIdleConns(),
		MaxOpenConns:              pgConf.GetMaxOpenConns(),
		Password:                  pgConf.GetPassword(),
		Port:                      pgConf.GetPort(),
		Username:                  pgConf.GetUser(),
		SSLMode:                   pgConf.GetSSLMode(),
		PingTimeout:               pgConf.GetPingTimeout(),
		SlowThreshold:             pgConf.GetLog().GetSlowThreshold(),
		LogLevel:                  pgConf.GetLog().GetLogLevel(),
		IgnoreRecordNotFoundError: pgConf.GetLog().GetIgnoreRecordNotFoundError(),
		ParameterizedQueries:      pgConf.GetLog().GetParameterizedQueries(),
		Colorful:                  pgConf.GetLog().GetColorful(),
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