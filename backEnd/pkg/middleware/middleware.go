package middleware

import (
	"context"
	"net/http"

	errorControllers "social-network/pkg/errorManagement/controllers"
	fileControllers "social-network/pkg/fileManagement/controllers"
	userModel "social-network/pkg/userManagement/models"
)

type contextKey string

const (
	CtxUserID       contextKey = "userID"
	CtxUserUUID     contextKey = "userUUID"
	CtxSessionID    contextKey = "sessionID"
	CtxProfileImage contextKey = "profileImage"
)

func CheckHttpRequest(checkFor, method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, userID, userUUID := checkSessionValidity(r)
		if checkFor == "guest" && sessionID != "" {
			errorControllers.CustomErrorHandler(w, r, "You are logged in", http.StatusBadRequest)
			return
		}
		if checkFor == "user" && userID == -1 {
			errorControllers.CustomErrorHandler(w, r, "Log in or register first", http.StatusBadRequest)
			return
		}
		if r.Method != method {
			errorControllers.ErrorHandler(w, r, errorControllers.MethodNotAllowedError)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, CtxUserID, userID)
		ctx = context.WithValue(ctx, CtxUserUUID, userUUID)
		ctx = context.WithValue(ctx, CtxSessionID, sessionID)

		next(w, r.WithContext(ctx))
	}
}

func checkSessionValidity(r *http.Request) (string, int, string) {
	sessionCookie, err := r.Cookie("session-id")
	if err != nil || sessionCookie == nil {
		return "", -1, ""
	}
	sessionID := sessionCookie.Value

	s, err := userModel.SelectActiveSessionBy("id", sessionID)
	if err != nil {
		return "", -1, ""
	}
	return sessionID, s.UserId, s.UserUUID
}

const maxUploadSize = 2 << 20 // 2 MB
func HandleProfileImageUpload(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

		err := r.ParseMultipartForm(maxUploadSize)
		if err != nil {
			errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
			return
		}

		profileImageFile, handler, err := r.FormFile("profile_image")
		if err != nil && err != http.ErrMissingFile {
			errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
			return
		}

		profileImage := ""
		if err == nil && handler.Size > 0 {
			defer profileImageFile.Close()

			if handler.Size > maxUploadSize {
				errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
				return
			}

			profileImage, err = fileControllers.FileUpload(profileImageFile, handler)
			if err != nil {
				errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
				return
			}
		}

		ctx := context.WithValue(r.Context(), CtxProfileImage, profileImage)
		next(w, r.WithContext(ctx))
	}
}

func GetSessionID(ctx context.Context) (string, bool) {
	sessionID, ok := ctx.Value(CtxSessionID).(string)
	return sessionID, ok
}

func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(CtxUserID).(int)
	return userID, ok
}

func GetUserUUID(ctx context.Context) (string, bool) {
	userUUID, ok := ctx.Value(CtxUserUUID).(string)
	return userUUID, ok
}
