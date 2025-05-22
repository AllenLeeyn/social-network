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
