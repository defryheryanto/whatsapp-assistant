package main

import (
	"github.com/defryheryanto/whatsapp-assistant/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupPostgresConnection() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.DatabaseConnectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
