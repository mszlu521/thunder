package database

import (
	"github.com/mszlu521/thunder/config"
	"github.com/mszlu521/thunder/db"
)

var (
	_db *db.MySQL
)

func InitDB(mysqlConf config.Mysql) {
	m := db.MySQL{
		Database:     mysqlConf.Database,
		Host:         mysqlConf.Host,
		MaxIdleConns: mysqlConf.MaxIdleConns,
		MaxOpenConns: mysqlConf.MaxOpenConns,
		Password:     mysqlConf.Password,
		Port:         mysqlConf.Port,
		Username:     mysqlConf.User,
		PingTimeout:  mysqlConf.PingTimeout,
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
