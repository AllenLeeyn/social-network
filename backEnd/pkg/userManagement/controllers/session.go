package controller

import (
	"net/http"
	"time"

	errorController "social-network/pkg/errorManagement/controllers"
	userModel "social-network/pkg/userManagement/models"
)

func ExpireSession(w http.ResponseWriter, sessionId string) error {
	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		Value:    "",              // Empty the cookie's value
		Expires:  time.Unix(0, 0), // Set expiration to a past date
		MaxAge:   -1,              // Invalidate the cookie immediately
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

	return userModel.UpdateSession(&session{
		IsActive:   false,
		ExpireTime: time.Now(),
		LastAccess: time.Now(),
		ID:         sessionId,
	})

	// need to close WS connection
	//m.CloseConn(s)
}

func generateSession(w http.ResponseWriter, r *http.Request, userId int) {
	session, err := userModel.InsertSession(&session{
		UserId:   userId,
		IsActive: true,
	})
	if err != nil {
		errorController.HandleErrorPage(w, r, errorController.InternalServerError)
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
