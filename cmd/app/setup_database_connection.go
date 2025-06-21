package main

import (
        "os"

        "gorm.io/driver/postgres"
        "gorm.io/gorm"
)

func setupPostgresConnection() (*gorm.DB, error) {
        dsn := os.Getenv("POSTGRES_DSN")
        if dsn == "" {
                dsn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
        }

        db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
