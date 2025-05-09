package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

func checkSessionValidity(w http.ResponseWriter, r *http.Request) (*http.Cookie, int) {
	sessionCookie, err := r.Cookie("session-id")
	if err != nil || sessionCookie == nil {
		log.Println(err)
		return nil, -1
	}
	sessionID := sessionCookie.Value
	s, err := db.SelectActiveSessionBy("id", sessionID)
	if err != nil || s.ExpireTime.Before(time.Now()) {
		expireSession(w, s)
		return nil, -1
	}
	return sessionCookie, s.UserID
}

func createSession(w http.ResponseWriter, user *user) {
	// generate a uuid for the session and set it into a cookie
	id, _ := uuid.NewV4()
	cookie := &http.Cookie{
		Name:     "session-id",
		Value:    id.String(),
		MaxAge:   7200,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	db.InsertSession(&session{
		ID:         id.String(),
		UserID:     user.ID,
		IsActive:   true,
		ExpireTime: time.Now().Add(2 * time.Hour),
	})
}

func extendSession(w http.ResponseWriter, sessionCookie *http.Cookie) {
	if sessionCookie == nil {
		return
	}
	// generate a uuid for the session and set it into a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		Value:    sessionCookie.Value,
		MaxAge:   7200,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	db.UpdateSession(&session{
		IsActive:   true,
		ExpireTime: time.Now().Add(2 * time.Hour),
		LastAccess: time.Now(),
		ID:         sessionCookie.Value,
	})
}

func expireSession(w http.ResponseWriter, s *session) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		Value:    "", // Empty the cookie's value
		MaxAge:   -1, // Invalidate the cookie immediately
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	db.UpdateSession(&session{
		IsActive:   false,
		ExpireTime: time.Now(),
		LastAccess: time.Now(),
		ID:         s.ID,
	})
	m.CloseConn(s)
}
