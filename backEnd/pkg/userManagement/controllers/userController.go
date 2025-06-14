package controller

import (
	"net/http"

	middleware "social-network/pkg/middleware"
	"social-network/pkg/utils"

	errorControllers "social-network/pkg/errorManagement/controllers"
	followingModel "social-network/pkg/followingManagement/models"
	userModel "social-network/pkg/userManagement/models"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type user = userModel.User

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	u := &user{}
	if err := utils.ReadJSON(r, u); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}

	if err := isValidRegistration(u); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	u.TypeId = 1

	password_hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	u.PasswordHash = string(password_hash)

	userId, userUUID, err := userModel.InsertUser(u)
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
	utils.ReturnJsonSuccess(w, "Registered successfully", userModel.UserView{
		UUID:         userUUID,
		NickName:     u.NickName,
		ProfileImage: u.ProfileImage,
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	u := &user{}
	if err := utils.ReadJSON(r, u); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
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

	if session, err := userModel.SelectActiveSessionBy("user_id", user.ID); err == nil && session != nil {
		expireSession(w, r, session)
	}
	generateSession(w, r, user.ID)

	utils.ReturnJsonSuccess(w, "login successfully", userModel.UserView{
		UUID:         user.UUID,
		NickName:     user.NickName,
		ProfileImage: user.ProfileImage,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionId, _ := middleware.GetSessionID(r.Context())
	session, err := userModel.SelectActiveSessionBy("id", sessionId)
	if err == nil && session != nil {
		expireSession(w, r, session)
	}
	utils.ReturnJsonSuccess(w, "Logged out successfully", nil)
}

func UserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	userID, isOk := middleware.GetUserID(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.UnauthorizedError)
		return
	}
	currentUserInfo, err := userModel.SelectUserByField("id", userID)
	if currentUserInfo == nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	u := &user{}
	if err := utils.ReadJSON(r, u); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
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

func ViewUserHandler(w http.ResponseWriter, r *http.Request) {
	userUUID, isOk := middleware.GetUserUUID(r.Context())
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.UnauthorizedError)
		return
	}

	tgtUUID, statusCode, err := GetTgtUUID(r, "api/user")

	if err != nil {
		tgtUser, err := userModel.SelectUser(tgtUUID)
		if err != nil {
			errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
			return
		}
		fStatus, _ := followingModel.SelectStatusByUUID(&followingModel.Following{
			LeaderUUID: tgtUUID, FollowerUUID: userUUID, GroupID: 0,
		})
		errorControllers.CustomErrorHandler(w, r,
			tgtUser.NickName+"|"+fStatus, statusCode)
		return
	}

	uProfile, err := userModel.SelectUser(tgtUUID)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}
	fStatus, _ := followingModel.SelectStatusByUUID(&followingModel.Following{
		LeaderUUID: tgtUUID, FollowerUUID: userUUID, GroupID: 0,
	})
	uProfile.Status = fStatus

	ExtendSession(w, r)
	utils.ReturnJsonSuccess(w, "user profile retrieved", uProfile)
}
