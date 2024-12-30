package database

import (
	"thunder/config"
	"thunder/db"
)

var (
	DB *db.MySQL
)

func InitDB() {
	m := db.MySQL{
		Database:     config.Conf.DB.Mysql.Database,
		Host:         config.Conf.DB.Mysql.Host,
		MaxIdleConns: config.Conf.DB.Mysql.MaxIdleConns,
		MaxOpenConns: config.Conf.DB.Mysql.MaxOpenConns,
		Password:     config.Conf.DB.Mysql.Password,
		Port:         config.Conf.DB.Mysql.Port,
		Username:     config.Conf.DB.Mysql.User,
		PingTimeout:  10,
	}
	err := m.Init()
	if err != nil {
		panic(err)
	}
	DB = &m
}
