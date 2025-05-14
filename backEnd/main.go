package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"social-network/pkg/db"
	"social-network/pkg/routes"

	userModel "social-network/pkg/userManagement/models"
)

var sqlDB *sql.DB

func init() {
	var err error
	sqlDB, err = db.OpenDB("sqlite3", "./pkg/db/social_network.db", "file://pkg/db/migrate")
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	userModel.NewUserModel(sqlDB)
}

func main() {
	routes.SetupRoutes(sqlDB)

	fmt.Println("Starting Forum on http://localhost:8080/...")
	log.Fatal(http.ListenAndServe(":8080", nil))
	sqlDB.Close()
}
