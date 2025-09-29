package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Postgres struct {
	GormConfig   *gorm.Config
	Username     string
	Password     string
	Host         string
	Port         int
	Database     string
	SSLMode      string
	PingTimeout  time.Duration
	MaxIdleConns int
	MaxOpenConns int
	GormDB       *gorm.DB
}

func (p *Postgres) Init() error {
	dsn := "host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai"
	dsn = fmt.Sprintf(dsn, p.Host, p.Username, p.Password, p.Database, p.Port, p.SSLMode)

	if p.GormConfig == nil {
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      true,        // Don't include params in the SQL logs
				Colorful:                  false,       // Disable color
			},
		)
		p.GormConfig = &gorm.Config{
			Logger: newLogger,
		}
	}
	db, err := gorm.Open(postgres.Open(dsn), p.GormConfig)
	if err != nil {
		return err
	}
	p.GormDB = db
	conn, err := db.DB()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), p.PingTimeout)
	defer cancel()
	err = conn.PingContext(ctx)
	if err != nil {
		return err
	}
	if p.MaxIdleConns != 0 {
		conn.SetMaxIdleConns(p.MaxIdleConns)
	}
	if p.MaxOpenConns != 0 {
		conn.SetMaxOpenConns(p.MaxOpenConns)
	}
	return nil
}