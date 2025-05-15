package controller

import (
	"net/http"

	errorControllers "social-network/pkg/errorManagement/controllers"
	middleware "social-network/pkg/middleware"
	userModel "social-network/pkg/userManagement/models"
	"social-network/pkg/utils"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type user = userModel.User

func (uc *UserController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	u := &user{}
	if err := utils.ReadJSON(w, r, u); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	if err := isValidUserInfo(u); err != nil {
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

	/* profileImageRaw := r.Context().Value("profileImage")
	profileImagePath, isOk := profileImageRaw.(string)
	if isOk && profileImagePath != "" {
		u.ProfileImage = sql.NullString{String: profileImagePath, Valid: true}
	} */

	userId, err := uc.um.InsertUser(u)
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

	uc.generateSession(w, r, userId)
	utils.ReturnJsonSuccess(w, "Registered successfully", nil)
}

func (uc *UserController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	u := &user{}
	if err := utils.ReadJSON(w, r, u); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	if err := isValidLogin(u); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	user := &user{}
	if u.Email != "" {
		user, _ = uc.um.SelectUserByField("email", u.Email)
	} else {
		user, _ = uc.um.SelectUserByField("nick_name", u.NickName)
	}

	if user == nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(u.Password)) != nil {
		errorControllers.CustomErrorHandler(w, r, "Incorrect username and/or password", http.StatusBadRequest)
		return
	}

	session, err := uc.um.SelectActiveSessionBy("user_id", user.ID)
	if err == nil {
		uc.ExpireSession(w, session.ID)
	}

	uc.generateSession(w, r, user.ID)
	utils.ReturnJsonSuccess(w, "Logged in successfully", nil)
}

func (uc *UserController) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionIdRaw := r.Context().Value(middleware.CtxSessionID)
	sessionId, isOk := sessionIdRaw.(string)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	if err := uc.ExpireSession(w, sessionId); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Logged out successfully", nil)
}

func (uc *UserController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	currentUserInfo, err := uc.um.SelectUserByField("id", userID)
	if currentUserInfo == nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	u := &user{}
	if err := utils.ReadJSON(w, r, u); err != nil {
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

	/* profileImageRaw := r.Context().Value("profileImage")
	profileImagePath, isOk := profileImageRaw.(string)
	if isOk && profileImagePath != "" {
		u.ProfileImage = sql.NullString{String: profileImagePath, Valid: true}
	} */

	updateError := uc.um.UpdateUser(currentUserInfo)
	if updateError != nil {
		errorControllers.CustomErrorHandler(w, r, updateError.Error(), http.StatusInternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Profile updated successfully", nil)
}
