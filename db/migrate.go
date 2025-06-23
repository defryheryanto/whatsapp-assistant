package main

import (
	"database/sql"
	"flag"
	"fmt"

	"github.com/defryheryanto/whatsapp-assistant/config"
	_ "github.com/golang-migrate/migrate/source"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattes/migrate/source/file"
)

func main() {
	downFlag := flag.Bool("down", false, "database migration down")
	flag.Parse()

	config.Init()
	fmt.Println("opening connection..")
	db, err := sql.Open("postgres", config.DatabaseConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	instance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("reading migration files")
	fSrc, err := (&file.File{}).Open("./db/migrations")
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithInstance("file", fSrc, "postgres", instance)
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
