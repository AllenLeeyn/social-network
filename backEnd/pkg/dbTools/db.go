package dbTools

import (
	"database/sql"
	"fmt"

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
func OpenDB(driver, dataSource string) (*DBContainer, error) {
	conn, err := sql.Open(driver, dataSource)
	if err != nil {
		return nil, err
	}
	conn.Exec("PRAGMA foreign_keys = ON;")
	return &DBContainer{conn: conn}, nil
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

// db.deleteAllusers() for testing purposes
func (db *DBContainer) DeleteAllUsers() error {
	query := "DELETE FROM users"
	_, err := db.conn.Exec(query)
	db.vacuumDB()
	return err
}

// db.deleteAllSessions() for testing purposes
func (db *DBContainer) DeleteAllSessions() error {
	query := "DELETE FROM sessions"
	_, err := db.conn.Exec(query)
	db.vacuumDB()
	return err
}

// db.deleteAllPosts() for testing purposes
func (db *DBContainer) DeleteAllComments() error {
	query := "DELETE FROM comments"
	_, err := db.conn.Exec(query)
	db.vacuumDB()
	return err
}

// db.deleteAllPosts() for testing purposes
func (db *DBContainer) DeleteAllPosts() error {
	query := "DELETE FROM posts"
	_, err := db.conn.Exec(query)
	db.vacuumDB()
	return err
}

// db.vacuumDB) for testing purposes
func (db *DBContainer) vacuumDB() error {
	query := "VACUUM"
	_, err := db.conn.Exec(query)
	return err
}
