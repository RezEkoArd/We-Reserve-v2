package config

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	DB *gorm.DB
}

func (cfg Config) ConnectDB() (*DB, error ) {
	// host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
	cfg.Psql.Host,
	cfg.Psql.User,
	cfg.Psql.Password,
	cfg.Psql.DBName,
	cfg.Psql.Port,
	cfg.Psql.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("[connection] 1 - cannot connect to database" + cfg.App.AppPort)
		return nil, err
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Error().Err(err).Msg("[connection] 2 - cannot connect to database")
		return nil,err
	}

	sqlDB.SetMaxOpenConns(int(cfg.Psql.DBMaxOpen))
	sqlDB.SetMaxIdleConns(int(cfg.Psql.DBMaxIdle))

	return &DB{DB: db}, nil
}