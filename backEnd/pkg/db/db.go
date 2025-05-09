package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/gofrs/uuid"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

// note: fields need to be validated before calling insert functions.
// The variables in the struct are initialized with default value,
// meaning they are not null when inserting to db.
// This means the variables/fields will not be recognise as empty/null by sql.

type DBContainer struct {
	conn       *sql.DB
	Categories []string // stores categories recorded in db.
}

// openDB() opens a sql database with the driver and dataSource given.
func OpenDB(driver, dataSource, migrateSource string) (*DBContainer, error) {
	if _, err := os.Stat(dataSource); os.IsNotExist(err) {
		file, err := os.Create(dataSource)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	conn, err := sql.Open(driver, dataSource)
	if err != nil {
		return nil, err
	}

	if err = migrateDB(conn, migrateSource); err != nil {
		log.Fatal("Error migrating database: ", err)
	}

	return &DBContainer{conn: conn}, nil
}

func migrateDB(conn *sql.DB, migrateSource string) error {
	conn.Exec("PRAGMA foreign_keys = OFF;")
	defer conn.Exec("PRAGMA foreign_keys = ON;")

	driver, err := sqlite.WithInstance(conn, &sqlite.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(migrateSource, "sqlite3", driver)
	if err != nil {
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
	if version == 0 || err == migrate.ErrNilVersion {
		if err = generateUuidTables(conn); err != nil {
			return err
		}
	}

	if err = m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No new migrations to apply.")
			err = nil
		}
	} else {
		version, _, _ = m.Version()
		log.Printf("Migrations applied successfully. New migration version: %d", version)
	}

	return err
}

func generateUuidTables(conn *sql.DB) error {
	log.Println("Initializing UUIDs for users and posts...")

	// Create the UUID tables if not already present
	createTables := `
	CREATE TABLE IF NOT EXISTS user_uuids (
		user_id INTEGER PRIMARY KEY,
		uuid TEXT NOT NULL UNIQUE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS post_uuids (
		post_id INTEGER PRIMARY KEY,
		uuid TEXT NOT NULL UNIQUE,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
	);`
	if _, err := conn.Exec(createTables); err != nil {
		return err
	}

	err := generateUuids(conn, "user")
	if err == nil {
		err = generateUuids(conn, "post")
	}
	return err
}

func generateUuids(conn *sql.DB, table string) error {
	// Insert UUIDs for users
	rows, err := conn.Query("SELECT id FROM " + table + "s")
	if err != nil {
		return checkErrNoRows(err)
	}
	uuidMap := make(map[int]string)

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		uuid, _ := uuid.NewV4()
		uuidMap[id] = uuid.String()
	}
	rows.Close()

	query := "INSERT INTO " + table + "_uuids (" + table + "_id, uuid) VALUES (?, ?)"
	for id, uuidStr := range uuidMap {
		if _, err = conn.Exec(query, id, uuidStr); err != nil {
			return err
		}
	}

	log.Println("UUID initialization complete.")
	return nil
}

// checkErrNoRows() checks if no result from sql query.
func checkErrNoRows(err error) error {
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (db *DBContainer) Close() {
	db.conn.Close()
}
