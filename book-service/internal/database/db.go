package database

import (
	"context"
	"fmt"
	"sync"

	"github.com/cenkalti/backoff"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"gitlab.com/narm-group/book-service/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var once sync.Once
var db *gorm.DB

func InitDB(cfg config.DBConfig) {
	once.Do(func() {
		var err error
		db, err = initDB(cfg)
		if err != nil {
			log.Fatalf("error db init: %v\n", err)
		}
	})
}

func GetDB(ctx context.Context) *gorm.DB {
	return db.WithContext(ctx)
}

func initDB(cfg config.DBConfig) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.User,
		cfg.Password,
	)

	log.Info("connecting to database")

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Warn("error getting sql db from gorm db")
		return
	}

	err = backoff.Retry(
		func() error {
			if err = sqlDB.Ping(); err != nil {
				log.Warnf("db ping failed : %v\n", err)
			}
			return err
		},
		backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3),
	)

	if err != nil {
		log.Error("db connection failed")
	} else {
		log.Info("connected to database")
	}

	return db, nil
}
