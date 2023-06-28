package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"regexp"

	_ "github.com/golang-migrate/migrate/source"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattes/migrate/source/file"
)

func main() {
	downFlag := flag.Bool("down", false, "database migration down")
	flag.Parse()

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/whatsapp_assistant.db", getAppRootDirectory()))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		panic(err)
	}

	fSrc, err := (&file.File{}).Open("./db/migrations")
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithInstance("file", fSrc, "sqlite3", instance)
	if err != nil {
		panic(err)
	}

	if *downFlag {
		fmt.Println("Rollback migration..")
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			panic(err)
		}
		version, _, _ := m.Version()
		fmt.Printf("Rollback complete to version %d.\n", version)
	} else {
		fmt.Println("Migrating migration..")
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			panic(err)
		}
		version, _, _ := m.Version()
		fmt.Printf("Migrate complete (version %d)\n", version)
	}
}

func getAppRootDirectory() string {
	projectName := regexp.MustCompile(`^(.*whatsapp-assistant)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	return string(rootPath)
}
