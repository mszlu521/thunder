package database

import (
	"github.com/mszlu521/thunder/config"
	"github.com/mszlu521/thunder/db"
)

var (
	_db *db.MySQL
)

func InitDB(mysqlConf *config.Mysql) {
	if mysqlConf == nil {
		return
	}

	m := db.MySQL{
		Database:     mysqlConf.GetDatabase(),
		Host:         mysqlConf.GetHost(),
		MaxIdleConns: mysqlConf.GetMaxIdleConns(),
		MaxOpenConns: mysqlConf.GetMaxOpenConns(),
		Password:     mysqlConf.GetPassword(),
		Port:         mysqlConf.GetPort(),
		Username:     mysqlConf.GetUser(),
		PingTimeout:  mysqlConf.GetPingTimeout(),
	}
	err := m.Init()
	if err != nil {
		panic(err)
	}
	_db = &m
}

func GetMysqlDB() *db.MySQL {
	return _db
}