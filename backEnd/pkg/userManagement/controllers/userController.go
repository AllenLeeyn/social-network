package controller

import (
	"net/http"

	errorControllers "social-network/pkg/errorManagement/controllers"
	followingModel "social-network/pkg/followingManagement/models"
	middleware "social-network/pkg/middleware"
	userModel "social-network/pkg/userManagement/models"
	"social-network/pkg/utils"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type user = userModel.User

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	u := &user{}
	if err := utils.ReadJSON(r, u); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	if err := isValidRegistration(u); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	u.TypeId = 1
	u.PasswordHash = string(password_hash)

	userId, err := userModel.InsertUser(u)
	if err != nil {
		if err.Error() == "email is already used" ||
			err.Error() == "nick name is already used" ||
			err.Error() == "email or nick name already exists" {
			errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
		} else {
			errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		}
		return
	}

	generateSession(w, r, userId)
	utils.ReturnJsonSuccess(w, "Registered successfully", nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	u := &user{}
	if err := utils.ReadJSON(r, u); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	if err := isValidLogin(u); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	user := &user{}
	if u.Email != "" {
		user, _ = userModel.SelectUserByField("email", u.Email)
	} else {
		user, _ = userModel.SelectUserByField("nick_name", u.NickName)
	}

	if user == nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(u.Password)) != nil {
		errorControllers.CustomErrorHandler(w, r, "Incorrect username and/or password", http.StatusBadRequest)
		return
	}

	session, err := userModel.SelectActiveSessionBy("user_id", user.ID)
	if err == nil && session != nil {
		ExpireSession(w, r, session)
	}

	generateSession(w, r, user.ID)
	utils.ReturnJsonSuccess(w, "Logged in successfully", nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionId, _ := middleware.GetSessionID(r.Context())
	session, err := userModel.SelectActiveSessionBy("id", sessionId)
	if err == nil && session != nil {
		ExpireSession(w, r, session)
	}
	utils.ReturnJsonSuccess(w, "Logged out successfully", nil)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	currentUserInfo, err := userModel.SelectUserByField("id", userID)
	if currentUserInfo == nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	u := &user{}
	if err := utils.ReadJSON(r, u); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := isValidUserInfo(u); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	currentUserInfo.FirstName = u.FirstName
	currentUserInfo.LastName = u.LastName
	currentUserInfo.NickName = u.NickName
	currentUserInfo.BirthDay = u.BirthDay
	currentUserInfo.AboutMe = u.AboutMe
	currentUserInfo.Visibility = u.Visibility

	updateError := userModel.UpdateUser(currentUserInfo)
	if updateError != nil {
		errorControllers.CustomErrorHandler(w, r, updateError.Error(), http.StatusInternalServerError)
		return
	}

	ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "Profile updated successfully", nil)
}

func ViewUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := userModel.SelectUsers()
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "user list retrieved", users)
}

func GetTgtUUID(r *http.Request) (string, int) {
	_, userID, userUUID, isOk := middleware.GetSessionCredentials(r.Context())
	if !isOk {
		return "", http.StatusInternalServerError
	}
	tgtUUID := r.URL.Query().Get("id")
	if tgtUUID == "" {
		tgtUUID = userUUID
	} else {
		if !userModel.IsPublic(tgtUUID) && !followingModel.IsFollower(userID, tgtUUID) {
			return "",
				http.StatusForbidden
		}
	}
	return tgtUUID, 200
}

func ViewUserHandler(w http.ResponseWriter, r *http.Request) {
	tgtUUID, statusCode := GetTgtUUID(r)
	if statusCode == http.StatusInternalServerError {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	} else if statusCode == http.StatusForbidden {
		errorControllers.CustomErrorHandler(w, r,
			"access denied: private profile and user is not follower",
			http.StatusForbidden)
		return
	}

	uProfile, err := userModel.SelectUser(tgtUUID)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "user profile retrieved", uProfile)
}
