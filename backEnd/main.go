package main

import (
	"fmt"
	"log"
	"net/http"
	"social-network/pkg/dbTools"
	"social-network/pkg/handlers"
	"social-network/pkg/messenger"
	"time"
)

var db *dbTools.DBContainer
var m messenger.Messenger

func init() {
	var err error
	db, err = dbTools.OpenDB("sqlite3", "./pkg/database/forum.db")
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	db.Categories, _ = db.SelectFieldFromTable("name", "categories")
	m = messenger.Start(db)
	handlers.Init(db, m)

	go expireSessionsTask()
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
	db.Close()
}

func expireSessionsTask() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sessions, err := db.SelectActiveSessions()
		if err != nil {
			fmt.Printf("Error selecting sessions: %v\n", err.Error())
		}
		for _, s := range *sessions {
			if time.Now().After(s.ExpireTime) {
				fmt.Printf("Expire session: %v\n", s.ID)

				s.IsActive = false
				err = db.UpdateSession(&s)
				if err != nil {
					break
				}
				err = m.CloseConn(&s)
				if err != nil {
					break
				}
			}
		}
		if err != nil {
			log.Printf("Error expiring sessions: %v", err)
		}
	}
}
