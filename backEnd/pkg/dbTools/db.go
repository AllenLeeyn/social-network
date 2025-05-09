package dbTools

import (
	"database/sql"
	"fmt"
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

/*
	DBContainer struct that comes with a set of functions.

This should be easier to reference the database and call its functions.
*/
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

	err = migrateDB(conn, migrateSource)
	if err != nil {
		log.Fatal("Error migrating database: ", err)
	}

	return &DBContainer{conn: conn}, nil
}

func migrateDB(conn *sql.DB, migrateSource string) error {
	conn.Exec("PRAGMA foreign_keys = OFF;")
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

	err = generateUuidTables(conn)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No new migrations to apply.")
		} else {
			return err
		}
	} else {
		version, _, _ = m.Version()
		log.Printf("Migrations applied successfully. New migration version: %d", version)
	}

	conn.Exec("PRAGMA foreign_keys = ON;")
	return nil
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
		_, err = conn.Exec(query, id, uuidStr)
		if err != nil {
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

// db.selectFieldFromTable() is a generic function to grab a column of data from a table.
func (db *DBContainer) SelectFieldFromTable(field, table string) ([]string, error) {
	rows, err := db.conn.Query("SELECT " + field + " FROM " + table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []string
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	if err := rows.Err(); err != nil {
		return nil, checkErrNoRows(err)
	}
	return values, nil
}

// db.isValidCategories*() check if given categories are valid with categories in db
func (db *DBContainer) isValidCategories(categories []int) error {
	if len(categories) == 0 {
		return fmt.Errorf("no categories")
	}
	for _, catID := range categories {
		if catID > len(db.Categories) || catID < 0 {
			return fmt.Errorf("invalid category")
		}
	}
	return nil
}

func (db *DBContainer) Close() {
	db.conn.Close()
}
