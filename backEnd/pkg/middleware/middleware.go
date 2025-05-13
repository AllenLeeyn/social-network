package middleware

import (
	"context"
	"net/http"
	"time"

	errorControllers "social-network/pkg/errorManagement/controllers"
	fileControllers "social-network/pkg/fileManagement/controllers"
	userContollers "social-network/pkg/userManagement/controllers"
	userModel "social-network/pkg/userManagement/models"
)

type contextKey string

const (
	ctxUserID              contextKey = "userID"
	ctxSessionID           contextKey = "sessionID"
	ContextKeyProfileImage contextKey = "profileImage"
)

func CheckHttpRequest(checkFor, method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, userID := checkSessionValidity(w, r)
		if checkFor == "guest" && sessionID != "" {
			errorControllers.HandleErrorPage(w, r, errorControllers.BadRequestError)
			return
		}
		if checkFor == "user" && userID == -1 {
			errorControllers.HandleErrorPage(w, r, errorControllers.UnauthorizedError)
			return
		}
		if r.Method != method {
			errorControllers.HandleErrorPage(w, r, errorControllers.MethodNotAllowedError)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxUserID, userID)
		ctx = context.WithValue(ctx, ctxSessionID, sessionID)

		next(w, r.WithContext(ctx))
	}
}

func checkSessionValidity(w http.ResponseWriter, r *http.Request) (string, int) {
	sessionCookie, err := r.Cookie("session-id")
	if err != nil || sessionCookie == nil {
		return "", -1
	}
	sessionID := sessionCookie.Value

	s, err := userModel.SelectActiveSessionBy("id", sessionID)
	if err != nil || s.ExpireTime.Before(time.Now()) {
		userContollers.ExpireSession(w, s.ID)
		return "", -1
	}
	return sessionID, s.UserId
}

const maxUploadSize = 2 << 20 // 2 MB
func HandleProfileImageUpload(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

		err := r.ParseMultipartForm(maxUploadSize)
		if err != nil {
			errorControllers.HandleErrorPage(w, r, errorControllers.BadRequestError)
			return
		}

		profileImageFile, handler, err := r.FormFile("profile_image")
		if err != nil && err != http.ErrMissingFile {
			errorControllers.HandleErrorPage(w, r, errorControllers.BadRequestError)
			return
		}

		profileImage := ""
		if err == nil && handler.Size > 0 {
			defer profileImageFile.Close()

			if handler.Size > maxUploadSize {
				errorControllers.HandleErrorPage(w, r, errorControllers.BadRequestError)
				return
			}

			profileImage, err = fileControllers.FileUpload(profileImageFile, handler)
			if err != nil {
				errorControllers.HandleErrorPage(w, r, errorControllers.InternalServerError)
				return
			}
		}

		ctx := context.WithValue(r.Context(), ContextKeyProfileImage, profileImage)
		next(w, r.WithContext(ctx))
	}
}
