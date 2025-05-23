package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	db "social-network/pkg/databaseManagement"
	"social-network/pkg/routes"

	chatModel "social-network/pkg/chatManagement/models"
	eventModel "social-network/pkg/eventManagement/models"
	followingModel "social-network/pkg/followingManagement/models"
	forumModel "social-network/pkg/forumManagement/models"
	groupModel "social-network/pkg/groupManagement/models"
	notificationModel "social-network/pkg/notificationManagement/models"
	userModel "social-network/pkg/userManagement/models"

	chatContollers "social-network/pkg/chatManagement/controllers"
	userControllers "social-network/pkg/userManagement/controllers"
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

	// need chat instance to closeConn when logOut
	userControllers.Initialize(cc)
}

func modelsInitDb(db *sql.DB) {
	chatModel.Initialize(sqlDB)
	eventModel.Initialize(db)
	followingModel.Initialize(sqlDB)
	forumModel.Initialize(db)
	groupModel.Initialize(sqlDB)
	notificationModel.Initialize(db)
	userModel.Initialize(db)
}

func main() {
	log.Println("\033[31mSetup routes\033[0m")
	routes.SetupRoutes(cc)

	fmt.Println("\033[32mStarting Forum on http://localhost:8080/...\033[0m")
	log.Fatal(http.ListenAndServe(":8080", nil))
	sqlDB.Close()
}
