package main

import (
	"fmt"
	"log"
	"net/http"
	"social-network/pkg/db"
	"social-network/pkg/handlers"
)

var dbConn *db.DBContainer

func init() {
	var err error
	dbConn, err = db.OpenDB("sqlite3", "./pkg/db/social_network.db", "file://pkg/db/migrate")
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
}

func main() {
	http.HandleFunc("/ws", handlers.WS)

	http.HandleFunc("/posts", handlers.Posts)
	http.HandleFunc("/post", handlers.Post)
	http.HandleFunc("/profile", handlers.Profile)

	http.HandleFunc("/signup", handlers.Signup)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/logout", handlers.LogOut)

	http.HandleFunc("/create-post", handlers.CreatePost)
	http.HandleFunc("/create-comment", handlers.CreateComment)
	http.HandleFunc("/feedback", handlers.CreateFeedback)

	fmt.Println("Starting Forum on http://localhost:8080/...")
	log.Fatal(http.ListenAndServe(":8080", nil))
	dbConn.Close()
}
