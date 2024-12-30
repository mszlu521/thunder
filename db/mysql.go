package db

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type MySQL struct {
	GormConfig   *gorm.Config
	Username     string
	Password     string
	Host         string
	Port         int
	Database     string
	PingTimeout  int
	MaxIdleConns int
	MaxOpenConns int
	GormDB       *gorm.DB
}

func (m *MySQL) Init() error {
	dsn := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn = fmt.Sprintf(dsn, m.Username, m.Password, m.Host, m.Port, m.Database)

	if m.GormConfig == nil {
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      true,        // Don't include params in the SQL log
				Colorful:                  false,       // Disable color
			},
		)
		m.GormConfig = &gorm.Config{
			Logger: newLogger,
		}
	}
	db, err := gorm.Open(mysql.Open(dsn), m.GormConfig)
	if err != nil {
		return err
	}
	m.GormDB = db
	conn, err := db.DB()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.PingTimeout)*time.Second)
	defer cancel()
	err = conn.PingContext(ctx)
	if err != nil {
		return err
	}
	if m.MaxIdleConns != 0 {
		conn.SetMaxIdleConns(m.MaxIdleConns)
	}
	if m.MaxOpenConns != 0 {
		conn.SetMaxOpenConns(m.MaxOpenConns)
	}
	return nil
}
