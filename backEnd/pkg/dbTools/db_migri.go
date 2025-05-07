package dbTools

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func (db *DBContainer) Migri() error {
	db.conn.Exec("PRAGMA foreign_keys = OFF;")
	driver, err := sqlite.WithInstance(db.conn, &sqlite.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://pkg/database/migri",
		"sqlite3",
		driver,
	)
	if err != nil {
		fmt.Println("here")
		return err
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return err
	}

	if dirty {
		log.Printf("Database is in a dirty state at version %d. Manual intervention may be required.", version)
		return err
	}

	log.Printf("Current migration version: %d", version)

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No new migrations to apply.")
		} else {
			return err
		}
	} else {
		log.Println("Migrations applied successfully.")
	}

	db.conn.Exec("PRAGMA foreign_keys = ON;")
	return nil
}
