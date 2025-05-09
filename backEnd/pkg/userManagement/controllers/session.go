package controller

import (
	"database/sql"
	errorManagementControllers "social-network/pkg/errorManagement/controllers"
	"time"

	"net/http"
	userManagementModels "social-network/pkg/userManagement/models"
)

func sessionGenerator(w http.ResponseWriter, r *http.Request, sqlDB *sql.DB, userId int) {
	session := &userManagementModels.Session{
		UserId:   userId,
		IsActive: true,
	}
	session, insertError := userManagementModels.InsertSession(sqlDB, session)
	if insertError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	// Set the session token in a cookie
	UserSetCookie(w, session.ID, session.ExpireTime)
}

// helper function to check for valid user session in cookie
func CheckLogin(w http.ResponseWriter, r *http.Request, sqlDB *sql.DB) (bool, userManagementModels.User, string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return false, userManagementModels.User{}, "", nil
	}

	sessionToken := cookie.Value
	user, expirationTime, selectError := userManagementModels.SelectSession(sqlDB, sessionToken)
	if selectError != nil {
		if selectError.Error() == "sql: no rows in result set" {
			deleteCookie(w, "session_token")
			return false, userManagementModels.User{}, "", nil
		} else {
			return false, userManagementModels.User{}, "", selectError
		}
	}

	// Check if the cookie has expired
	if time.Now().After(expirationTime) {
		// Cookie expired, redirect to login
		return false, userManagementModels.User{}, "", nil
	}

	return true, user, sessionToken, nil
}

func deleteCookie(w http.ResponseWriter, cookieName string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",              // Optional but recommended
		Expires:  time.Unix(0, 0), // Set expiration to a past date
		MaxAge:   -1,              // Ensure immediate removal
		Path:     "/",             // Must match the original cookie path
		HttpOnly: true,
		Secure:   false,
	})
}

func UserSetCookie(w http.ResponseWriter, sessionToken string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})
}
