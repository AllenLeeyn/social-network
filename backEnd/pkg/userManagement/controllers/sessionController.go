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

func ExpireSession(w http.ResponseWriter, r *http.Request, sessionId string) error {
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
		ID:         sessionId,
	})
	if err != nil {
		return err
	}

	userUUId, isOk := middleware.GetUserUUID(r.Context())
	if !isOk {
		return fmt.Errorf("error getting userId from context")
	}
	return chatController.CloseConn(userUUId)
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
