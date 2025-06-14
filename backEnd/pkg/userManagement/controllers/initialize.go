package controller

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	middleware "social-network/pkg/middleware"
	"social-network/pkg/utils"

	chatContollers "social-network/pkg/chatManagement/controllers"
	followingModel "social-network/pkg/followingManagement/models"
	userModel "social-network/pkg/userManagement/models"
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

func GetTgtUUID(r *http.Request, basePath string) (string, int, error) {
	_, userID, userUUID, isOk := middleware.GetSessionCredentials(r.Context())
	if !isOk {
		return "", http.StatusInternalServerError, fmt.Errorf("internal server error")
	}
	tgtUUID, err := utils.ExtractUUIDFromUrl(r.URL.Path, basePath)
	if err != nil {
		return "", http.StatusNotFound, fmt.Errorf("page not found")
	}
	if tgtUUID == "" || tgtUUID == userUUID {
		tgtUUID = userUUID

	} else {
		if !userModel.IsPublic(tgtUUID) && !followingModel.IsFollower(userID, tgtUUID) {
			return tgtUUID, http.StatusForbidden,
				fmt.Errorf("access denied: private profile and user is not follower")
		}
	}
	return tgtUUID, http.StatusOK, nil
}
