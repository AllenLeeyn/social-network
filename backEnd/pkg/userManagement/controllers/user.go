package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	errorControllers "social-network/pkg/errorManagement/controllers"
	userModel "social-network/pkg/userManagement/models"
	"social-network/pkg/utils"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type user = userModel.User
type session = userModel.Session

func isValidUserInfo(u *user) error {
	isValid := false

	if u.FirstName, isValid = utils.IsValidUseName(u.FirstName); !isValid {
		return errors.New("First name must be between 3 to 16 alphanumeric characters, '_' or '-'")
	}
	if u.LastName, isValid = utils.IsValidUseName(u.LastName); !isValid {
		return errors.New("First name must be between 3 to 16 alphanumeric characters, '_' or '-'")
	}
	if u.NickName, isValid = utils.IsValidUseName(u.NickName); !isValid && u.NickName != "" {
		return errors.New("Nick name must be between 3 to 16 alphanumeric characters, '_' or '-'")
	}
	if u.Password, isValid = utils.IsValidPsswrd(u.Password); !isValid {
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
		return errors.New("About me is limited to 500 characters.")
	}

	return nil
}

func isValidRegistration(u *user) error {
	if err := isValidUserInfo(u); err != nil {
		return err
	}
	if user, _ := userModel.SelectUserByField("email", u.Email); user != nil {
		return errors.New("email is already used")
	}
	if user, _ := userModel.SelectUserByField("nick_name", u.NickName); user != nil {
		return errors.New("name is already used")
	}

	return nil
}

// need to handle image uploaded on registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	u := &user{}
	if err := utils.GetJSON(w, r, u); err != nil {
		errorControllers.HandleErrorPage(w, r, errorControllers.InternalServerError)
		return
	}

	if err := isValidRegistration(u); err != nil {
		utils.ReturnJson(w, utils.Result{
			Success:    false,
			Message:    err.Error(),
			HttpStatus: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		errorControllers.HandleErrorPage(w, r, errorControllers.InternalServerError)
		return
	}

	u.TypeId = 1
	u.PasswordHash = string(password_hash)

	profileImageRaw := r.Context().Value("profileImage")
	profileImagePath, isOk := profileImageRaw.(string)
	if isOk && profileImagePath != "" {
		u.ProfileImage = sql.NullString{String: profileImagePath, Valid: true}
	}

	// Insert a record while checking duplicates
	userId, err := userModel.InsertUser(u)
	if err != nil {
		fmt.Println(err.Error())
		errorControllers.HandleErrorPage(w, r, errorControllers.InternalServerError)
		return
	}

	generateSession(w, r, userId)

	utils.ReturnJson(w, utils.Result{
		Success: true,
		Message: "Logged in successfully",
		Data:    nil,
	})
	return
}

func isValidLogin(u *user) error {
	isValid := false

	if u.Email != "" {
		if u.Email, isValid = utils.IsValidEmail(u.Email); !isValid {
			return errors.New("invalid email")
		}
	} else if u.NickName != "" {
		if u.NickName, isValid = utils.IsValidUseName(u.NickName); !isValid {
			return errors.New("invalid user name")
		}
	} else {
		return errors.New("Email or user name is required")
	}

	if u.Password, isValid = utils.IsValidPsswrd(u.Password); !isValid {
		return errors.New("password must be 8 characters or longer.\n" +
			"Include at least a lower case character, an upper case character, a number and one of '@$!%*?&'")
	}
	return nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	u := &user{}
	if err := utils.GetJSON(w, r, u); err != nil {
		errorControllers.HandleErrorPage(w, r, errorControllers.InternalServerError)
		return
	}

	if err := isValidLogin(u); err != nil {
		errorControllers.HandleErrorPage(w, r, errorControllers.BadRequestError)
		return
	}

	user := &user{}
	if u.Email != "" {
		user, _ = userModel.SelectUserByField("email", u.Email)
	} else {
		user, _ = userModel.SelectUserByField("nick_name", u.NickName)
	}

	if user == nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(u.Password)) != nil {
		utils.ReturnJson(w, utils.Result{
			Success:    false,
			Message:    "Incorrect username and/or password",
			HttpStatus: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	session, err := userModel.SelectActiveSessionBy("user_id", user.ID)
	if err == nil {
		ExpireSession(w, session.ID)
	}

	generateSession(w, r, user.ID)
	utils.ReturnJson(w, utils.Result{
		Success: true,
		Message: "Logged in successfully",
		Data:    nil,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	sessionIdRaw := r.Context().Value("sessionID")
	sessionId, isOk := sessionIdRaw.(string)
	if !isOk {
		errorControllers.HandleErrorPage(w, r, errorControllers.InternalServerError)
		return
	}
	if err := ExpireSession(w, sessionId); err != nil {
		errorControllers.HandleErrorPage(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJson(w, utils.Result{
		Success: true,
		Message: "Logged out successfully",
		Data:    nil,
	})
	return
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("userID")
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.HandleErrorPage(w, r, errorControllers.InternalServerError)
		return
	}

	curUserInfo, _ := userModel.SelectUserByField("id", userID)
	if curUserInfo == nil {
		errorControllers.HandleErrorPage(w, r, errorControllers.InternalServerError)
		return
	}

	u := &user{}
	if err := utils.GetJSON(w, r, u); err != nil {
		errorControllers.HandleErrorPage(w, r, errorControllers.InternalServerError)
		return
	}

	if err := isValidUserInfo(u); err != nil {
		utils.ReturnJson(w, utils.Result{
			Success:    false,
			Message:    err.Error(),
			HttpStatus: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}

	curUserInfo.FirstName = u.FirstName
	curUserInfo.LastName = u.LastName
	curUserInfo.NickName = u.NickName
	curUserInfo.BirthDay = u.BirthDay
	curUserInfo.AboutMe = u.AboutMe
	curUserInfo.Visibility = u.Visibility

	profileImageRaw := r.Context().Value("profileImage")
	profileImagePath, isOk := profileImageRaw.(string)
	if isOk && profileImagePath != "" {
		curUserInfo.ProfileImage = sql.NullString{String: profileImagePath, Valid: true}
	}

	// Update a record while checking duplicates
	updateError := userModel.UpdateUser(curUserInfo)
	if updateError != nil {
		errorControllers.HandleErrorPage(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJson(w, utils.Result{
		Success: true,
		Message: "Profile updated successfully",
		Data:    nil,
	})
	return
}
