package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"social-network/pkg/db"
	"social-network/pkg/routes"

	chatContollers "social-network/pkg/chatManagement/controllers"
	chatModel "social-network/pkg/chatManagement/models"
	userControllers "social-network/pkg/userManagement/controllers"
	userModel "social-network/pkg/userManagement/models"
)

var sqlDB *sql.DB

func init() {
	log.Println("\033[31mInitialise database\033[0m")
	var err error
	sqlDB, err = db.OpenDB("sqlite3", "./pkg/db/social_network.db", "file://pkg/db/migrate")
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	log.Println("\033[31mInitialise models\033[0m")
	userModel.Initialize(sqlDB)
	chatModel.Initialize(sqlDB)
}

func main() {

	log.Println("\033[31mInitialise controllers\033[0m")
	cc := chatContollers.Initialize()
	userControllers.Initialize(cc)

	log.Println("\033[31mSetup routes\033[0m")
	routes.SetupRoutes(sqlDB, cc)

	fmt.Println("\033[32mStarting Forum on http://localhost:8080/...\033[0m")
	log.Fatal(http.ListenAndServe(":8080", nil))
	sqlDB.Close()
}
