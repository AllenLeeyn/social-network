package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	userModel "social-network/pkg/userManagement/models"
	"social-network/pkg/utils"
)

type UserController struct {
	um *userModel.UserModel
}

func NewUserController(dbMain *sql.DB) *UserController {
	return &UserController{
		um: userModel.NewUserModel(dbMain),
	}
}

func isValidUserInfo(u *user) error {
	isValid := false

	if u.FirstName, isValid = utils.IsValidUserName(u.FirstName); !isValid {
		return errors.New("first name must be between 3 to 16 alphanumeric characters, '_' or '-'")
	}
	if u.LastName, isValid = utils.IsValidUserName(u.LastName); !isValid {
		return errors.New("last name must be between 3 to 16 alphanumeric characters, '_' or '-'")
	}
	if u.NickName != "" {
		u.NickName = generateNickName(u)
	}
	if u.NickName, isValid = utils.IsValidUserName(u.NickName); !isValid {
		return errors.New("nick name must be between 3 to 16 alphanumeric characters, '_' or '-'")
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

func isValidLogin(u *user) error {
	isValid := false

	if u.Email != "" {
		if u.Email, isValid = utils.IsValidEmail(u.Email); !isValid {
			return errors.New("invalid email")
		}
	} else if u.NickName != "" {
		if u.NickName, isValid = utils.IsValidUserName(u.NickName); !isValid {
			return errors.New("nick name must be between 3 to 16 alphanumeric characters, '_' or '-'")
		}
	} else {
		return errors.New("email or user name is required")
	}

	if u.Password, isValid = utils.IsValidPassword(u.Password); !isValid {
		return errors.New("password must be 8 characters or longer.\n" +
			"Include at least a lower case character, an upper case character, a number and one of '@$!%*?&'")
	}
	return nil
}

func generateNickName(u *user) string {
	base := strings.ToLower(u.FirstName)
	if len(base) > 10 {
		base = base[:10]
	}

	initial := ""
	if len(u.LastName) > 0 {
		initial = strings.ToUpper(string(u.LastName[0]))
	}

	num := rand.Intn(10000)
	return fmt.Sprintf("%s%s_%d", base, initial, num)
}
