package controller

import (
	"fmt"
	"net/http"
	"time"

	middleware "social-network/pkg/middleware"

	errorController "social-network/pkg/errorManagement/controllers"
	userModel "social-network/pkg/userManagement/models"
)

type session = userModel.Session

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
		Path:     "/",
		//Secure:   false,
		//SameSite: http.SameSiteNoneMode,
	})
}

func ExtendSession(w http.ResponseWriter, r *http.Request) error {
	sessionID, _ := middleware.GetSessionID(r.Context())
	// generate a uuid for the session and set it into a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		Value:    sessionID,
		Expires:  time.Now().Add(time.Hour * 24 * 365 * 10),
		MaxAge:   315360000, // 10 years
		Path:     "/",
		HttpOnly: true,
		//Secure:   false,
		//SameSite: http.SameSiteNoneMode,
	})
	return userModel.UpdateSession(&session{
		IsActive:   true,
		ExpireTime: time.Now().Add(time.Hour * 24 * 365 * 10),
		LastAccess: time.Now(),
		ID:         sessionID,
	})
}

func expireSession(w http.ResponseWriter, r *http.Request, s *session) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		Value:    "",         // Empty the cookie's value
		Expires:  time.Now(), // Set expiration to a past date
		MaxAge:   -1,         // Invalidate the cookie immediately
		HttpOnly: true,
		Path:     "/",
		//Secure:   false,
		//SameSite: http.SameSiteNoneMode,
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
