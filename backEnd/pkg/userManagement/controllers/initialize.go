package controller

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	chatContollers "social-network/pkg/chatManagement/controllers"
	followingModel "social-network/pkg/followingManagement/models"
	middleware "social-network/pkg/middleware"
	userModel "social-network/pkg/userManagement/models"
	"social-network/pkg/utils"
)

var chatController *chatContollers.ChatController

func Initialize(cc *chatContollers.ChatController) {
	log.Println("\033[35mInitlise user controller\033[0m")
	chatController = cc
}

func isValidUserInfo(u *user) error {
	isValid := false

	if u.FirstName, isValid = utils.IsValidUserName(u.FirstName); !isValid {
		return errors.New("first name must be between 3 to 16 alphanumeric characters, '_' or '-'")
	}
	if u.LastName, isValid = utils.IsValidUserName(u.LastName); !isValid {
		return errors.New("last name must be between 3 to 16 alphanumeric characters, '_' or '-'")
	}
	if u.NickName == "" {
		u.NickName = generateNickName(u)
	}
	if u.NickName, isValid = utils.IsValidUserName(u.NickName); !isValid {
		return errors.New("nick name must be between 3 to 16 alphanumeric characters, '_' or '-'")
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

func isPasswordConfirmed(u *user) error {
	isValid := false
	if u.Password, isValid = utils.IsValidPassword(u.Password); !isValid {
		return errors.New("password must be 8 characters or longer.\n" +
			"Include at least a lower case character, an upper case character, a number and one of '@$!%*?&'")
	}
	if u.Password != u.ConfirmPassword {
		return errors.New("passwords do not match")
	}
	return nil
}

func isValidRegistration(u *user) error {
	if err := isValidUserInfo(u); err != nil {
		return err
	}
	if err := isPasswordConfirmed(u); err != nil {
		return err
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

func GetTgtUUID(r *http.Request, basePath string) (string, int) {
	_, userID, userUUID, isOk := middleware.GetSessionCredentials(r.Context())
	if !isOk {
		return "internal server error", http.StatusInternalServerError
	}
	tgtUUID, err := utils.ExtractUUIDFromUrl(r.URL.Path, basePath)
	if err != nil {
		return "page not found", http.StatusNotFound
	}
	if tgtUUID == "" {
		tgtUUID = userUUID

	} else {
		if !userModel.IsPublic(tgtUUID) && !followingModel.IsFollower(userID, tgtUUID) {
			return "access denied: private profile and user is not follower", http.StatusForbidden
		}
	}
	return tgtUUID, http.StatusOK
}
