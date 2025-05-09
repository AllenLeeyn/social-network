package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"social-network/pkg/db"
	"social-network/pkg/routes"
)

var sqlDB *sql.DB

func init() {
	var err error
	sqlDB, err = db.OpenDB("sqlite3", "./pkg/db/social_network.db", "file://pkg/db/migrate")
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
}

func main() {
	/*
		http.HandleFunc("/posts", handlers.Posts)
		http.HandleFunc("/post", handlers.Post)
		http.HandleFunc("/profile", handlers.Profile)

		http.HandleFunc("/signup", handlers.Signup)
		http.HandleFunc("/login", handlers.Login)
		http.HandleFunc("/logout", handlers.LogOut)

		http.HandleFunc("/create-post", handlers.CreatePost)
		http.HandleFunc("/create-comment", handlers.CreateComment)
		http.HandleFunc("/feedback", handlers.CreateFeedback)
	*/
	routes.SetupRoutes(sqlDB)

	fmt.Println("Starting Forum on http://localhost:8080/...")
	log.Fatal(http.ListenAndServe(":8080", nil))
	sqlDB.Close()
}
