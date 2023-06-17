package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupSQLiteConnection(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path))
	if err != nil {
		return nil, err
	}

	return db, nil
}
