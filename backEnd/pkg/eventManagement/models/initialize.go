package models

import (
	"database/sql"
	"log"
)

var sqlDB *sql.DB

func Initialize(dbMain *sql.DB) {
	log.Println("\033[34mInitlise following model\033[0m")
	sqlDB = dbMain
}

// checkErrNoRows() checks if no result from sql query.
func checkErrNoRows(err error) error {
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}
