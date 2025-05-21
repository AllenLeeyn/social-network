package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	errorControllers "social-network/pkg/errorManagement/controllers"
	middleware "social-network/pkg/middleware"
	userModel "social-network/pkg/userManagement/models"
	"social-network/pkg/utils"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type user = userModel.User

func isValidUserInfo(u *user) error {
	isValid := false

	if u.FirstName, isValid = utils.IsValidUserName(u.FirstName); !isValid {
		return errors.New("first name must be between 3 to 16 alphanumeric characters, '_' or '-'")
	}
	if u.LastName, isValid = utils.IsValidUserName(u.LastName); !isValid {
		return errors.New("last name must be between 3 to 16 alphanumeric characters, '_' or '-'")
	}
	if u.NickNameForm, isValid = utils.IsValidUserName(u.NickNameForm); !isValid && u.NickNameForm != "" {
		return errors.New("nick name must be between 3 to 16 alphanumeric characters, '_' or '-'")
	}
	if u.NickNameForm != "" {
		u.NickName = sql.NullString{Valid: true, String: u.NickNameForm}
	}
	if u.Password, isValid = utils.IsValidPassword(u.Password); !isValid {
		return errors.New("password must be 8 characters or longer.\n" +
			"Include at least a lower case character, an upper case character, a number and one of '@$!%*?&'")
	}
	if u.Password != u.ConfirmPassword {
		return errors.New("passwords do not match")
	}
	if u.Email, isValid = utils.IsValidEmail(u.Email); !isValid {
		return errors.New("invalid email")
	}
	if u.Visibility != "public" {
		u.Visibility = "private"
	}
	if u.Gender != "Male" && u.Gender != "Female" {
		u.Visibility = "Other"
	}
	if u.AboutMe, isValid = utils.IsValidContent(u.AboutMe, 0, 500); !isValid {
		return errors.New("about me is limited to 500 characters")
	}

	return nil
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	u := &user{}
	if err := utils.ReadJSON(w, r, u); err != nil {
		fmt.Println(err)
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

func isValidLogin(u *user) error {
	isValid := false

	if u.Email != "" {
		if u.Email, isValid = utils.IsValidEmail(u.Email); !isValid {
			return errors.New("invalid email")
		}
	} else if u.NickNameForm != "" {
		if u.NickNameForm, isValid = utils.IsValidUserName(u.NickNameForm); !isValid {
			return errors.New("nick name must be between 3 to 16 alphanumeric characters, '_' or '-'")
		}
		u.NickName = sql.NullString{Valid: true, String: u.NickNameForm}
	} else {
		return errors.New("email or user name is required")
	}

	if u.Password, isValid = utils.IsValidPassword(u.Password); !isValid {
		return errors.New("password must be 8 characters or longer.\n" +
			"Include at least a lower case character, an upper case character, a number and one of '@$!%*?&'")
	}
	return nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	u := &user{}
	if err := utils.ReadJSON(w, r, u); err != nil {
		fmt.Println(err)
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
		user, _ = userModel.SelectUserByField("nick_name", u.NickName.String)
	}

	if user == nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(u.Password)) != nil {
		errorControllers.CustomErrorHandler(w, r, "Incorrect username and/or password", http.StatusBadRequest)
		return
	}

	session, err := userModel.SelectActiveSessionBy("user_id", user.ID)
	if err == nil {
		ExpireSession(w, session.ID)
	}

	generateSession(w, r, user.ID)
	utils.ReturnJsonSuccess(w, "Logged in successfully", nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionIdRaw := r.Context().Value(middleware.CtxSessionID)
	sessionId, isOk := sessionIdRaw.(string)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	if err := ExpireSession(w, sessionId); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Logged out successfully", nil)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
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

	updateError := userModel.UpdateUser(currentUserInfo)
	if updateError != nil {
		errorControllers.CustomErrorHandler(w, r, updateError.Error(), http.StatusInternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Profile updated successfully", nil)
}
