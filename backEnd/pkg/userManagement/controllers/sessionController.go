package controller

import (
	"fmt"
	"net/http"
	"time"

	errorController "social-network/pkg/errorManagement/controllers"
	middleware "social-network/pkg/middleware"
	userModel "social-network/pkg/userManagement/models"
)

type session = userModel.Session

func ExtendSession(w http.ResponseWriter, r *http.Request) error {

	sessionID, _ := middleware.GetSessionID(r.Context())
	// generate a uuid for the session and set it into a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		Value:    sessionID,
		Expires:  time.Now().Add(2 * time.Hour),
		MaxAge:   7200,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	return userModel.UpdateSession(&session{
		IsActive:   true,
		ExpireTime: time.Now().Add(2 * time.Hour),
		LastAccess: time.Now(),
		ID:         sessionID,
	})
}

func ExpireSession(w http.ResponseWriter, r *http.Request, s *session) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		Value:    "",         // Empty the cookie's value
		Expires:  time.Now(), // Set expiration to a past date
		MaxAge:   -1,         // Invalidate the cookie immediately
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

	err := userModel.UpdateSession(&session{
		IsActive:   false,
		ExpireTime: time.Now(),
		LastAccess: time.Now(),
		ID:         s.ID,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(chatController.CloseConn(s.UserUUID))
}

func generateSession(w http.ResponseWriter, r *http.Request, userId int) {
	session, err := userModel.InsertSession(&session{
		UserId:   userId,
		IsActive: true,
	})
	if err != nil {
		errorController.ErrorHandler(w, r, errorController.InternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		Value:    session.ID,
		Expires:  session.ExpireTime,
		MaxAge:   int(time.Until(session.ExpireTime).Seconds()),
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})
}
