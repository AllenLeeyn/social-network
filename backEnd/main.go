package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	chatContollers "social-network/pkg/chatManagement/controllers"
	chatModel "social-network/pkg/chatManagement/models"
	db "social-network/pkg/databaseManagement"
	followingModel "social-network/pkg/followingManagement/models"
	groupModel "social-network/pkg/groupManagement/models"
	"social-network/pkg/routes"
	userControllers "social-network/pkg/userManagement/controllers"

	categoryModel "social-network/pkg/forumManagement/models"
	notificationModel "social-network/pkg/notificationManagement/models"
	userModel "social-network/pkg/userManagement/models"
)

var sqlDB *sql.DB
var cc = chatContollers.Initialize()

func init() {
	log.Println("\033[31mInitialise database\033[0m")
	var err error
	sqlDB, err = db.OpenDB("sqlite3",
		"./pkg/databaseManagement/social_network.db",
		"file://pkg/databaseManagement/migrate")
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	log.Println("\033[31mInitialise models\033[0m")
	modelsInitDb(sqlDB)

	log.Println("\033[31mInitialise controllers\033[0m")
	userControllers.Initialize(cc)
}

func modelsInitDb(db *sql.DB) {
	userModel.Initialize(db)
	chatModel.Initialize(sqlDB)
	groupModel.Initialize(sqlDB)
	followingModel.Initialize(sqlDB)
	categoryModel.Initialize(db)
	notificationModel.Initialize(db)
}

func main() {
	log.Println("\033[31mSetup routes\033[0m")
	routes.SetupRoutes(cc)

	fmt.Println("\033[32mStarting Forum on http://localhost:8080/...\033[0m")
	log.Fatal(http.ListenAndServe(":8080", nil))
	sqlDB.Close()
}
