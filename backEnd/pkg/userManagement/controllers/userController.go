package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	errorManagementControllers "social-network/pkg/errorManagement/controllers"
	fileManagementControllers "social-network/pkg/fileManagement/controllers"
	"sync"

	"net/http"
	userManagementModels "social-network/pkg/userManagement/models"
	"social-network/pkg/utils"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var OnlineUsers = make(map[*websocket.Conn]string) // Map of online users (connected to WS) to usernames
var Mutex = &sync.Mutex{}                          // Mutex to handle concurrent access to OnlineUsers

type AuthPageErrorData struct {
	ErrorMessage string
}

func RegisterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, _, _, checkLoginError := CheckLogin(w, r, db)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		res := utils.Result{
			Success: true,
			Message: "You are logged in",
			Data:    nil,
		}
		utils.ReturnJson(w, res)
		return
	}

	err := r.ParseMultipartForm(0)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	nick_name := utils.SanitizeInput(r.FormValue("nick_name"))
	first_name := utils.SanitizeInput(r.FormValue("first_name"))
	last_name := utils.SanitizeInput(r.FormValue("last_name"))
	gender := utils.SanitizeInput(r.FormValue("gender"))
	birthday_str := utils.SanitizeInput(r.FormValue("birthday"))
	email := utils.SanitizeInput(r.FormValue("email"))
	password := utils.SanitizeInput(r.FormValue("password"))
	about_me := utils.SanitizeInput(r.FormValue("about_me"))
	visibility := utils.SanitizeInput(r.FormValue("visibility"))
	if len(nick_name) == 0 || len(first_name) == 0 || len(last_name) == 0 || len(gender) == 0 || len(birthday_str) == 0 || len(email) == 0 || len(password) == 0 {
		res := utils.Result{
			Success:    false,
			Message:    "nick_name, first_name, last_name, gender, birthday, email and password are required.",
			HttpStatus: http.StatusOK,
			Data:       nil,
		}
		utils.ReturnJson(w, res)
		return
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		res := utils.Result{
			Success:    false,
			Message:    "Invalid email address!",
			HttpStatus: http.StatusOK,
			Data:       nil,
		}
		utils.ReturnJson(w, res)
		return
	}
	if len(visibility) == 0 {
		visibility = "private"
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	birthday, err := time.Parse("2006-01-02", birthday_str) // Adjust format as needed
	if err != nil {
		res := utils.Result{
			Success:    false,
			Message:    "Invalid birth date format. Use YYYY-MM-DD.",
			HttpStatus: http.StatusOK,
			Data:       nil,
		}
		utils.ReturnJson(w, res)
		return
	}

	newUser := &userManagementModels.User{
		NickName:     nick_name,
		FirstName:    first_name,
		LastName:     last_name,
		Gender:       gender,
		BirthDay:     birthday,
		Email:        email,
		PasswordHash: string(password_hash),
		AboutMe:      about_me,
		Visibility:   visibility,
	}

	// Insert a record while checking duplicates
	userId, insertError := userManagementModels.InsertUser(db, newUser)
	if insertError != nil {
		if insertError.Error() == "duplicateEmail" {
			res := utils.Result{
				Success:    false,
				Message:    "User with this email already exists!",
				HttpStatus: http.StatusOK,
				Data:       nil,
			}
			utils.ReturnJson(w, res)
			return
		} else if insertError.Error() == "duplicateNickName" {
			res := utils.Result{
				Success:    false,
				Message:    "User with this username already exists!",
				HttpStatus: http.StatusOK,
				Data:       nil,
			}
			utils.ReturnJson(w, res)
			return
		} else {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		}
		return
	}

	sessionGenerator(w, r, db, userId)

	res := utils.Result{
		Success: true,
		Message: "Logged in successfully",
		Data:    nil,
	}
	utils.ReturnJson(w, res)
	return
}

func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, _, _, checkLoginError := CheckLogin(w, r, db)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if loginStatus {
		res := utils.Result{
			Success: true,
			Message: "You are logged in",
			Data:    nil,
		}
		utils.ReturnJson(w, res)
		return
	}

	err := r.ParseMultipartForm(0)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	nick_name := utils.SanitizeInput(r.FormValue("nick_name"))
	password := utils.SanitizeInput(r.FormValue("password"))
	if len(nick_name) == 0 || len(password) == 0 {
		res := utils.Result{
			Success:    false,
			Message:    "NickName and password are required.",
			HttpStatus: http.StatusOK,
			Data:       nil,
		}
		utils.ReturnJson(w, res)
		return
	}

	// Insert a record while checking duplicates
	authStatus, userId, authError := userManagementModels.AuthenticateUser(db, nick_name, password)
	if authError != nil {
		res := utils.Result{
			Success:    false,
			Message:    authError.Error(),
			HttpStatus: http.StatusOK,
			Data:       nil,
		}
		utils.ReturnJson(w, res)
		return
	}
	if authStatus {
		sessionGenerator(w, r, db, userId)
	}

	res := utils.Result{
		Success: true,
		Message: "Logged in successfully",
		Data:    nil,
	}
	utils.ReturnJson(w, res)
	return
}

func sessionGenerator(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	session := &userManagementModels.Session{
		UserId:   userId,
		IsActive: true,
	}
	session, insertError := userManagementModels.InsertSession(db, session)
	if insertError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	// Set the session token in a cookie
	UserSetCookie(w, session.ID, session.ExpireTime)
}

// helper function to check for valid user session in cookie
func CheckLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) (bool, userManagementModels.User, string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return false, userManagementModels.User{}, "", nil
	}

	sessionToken := cookie.Value
	user, expirationTime, selectError := userManagementModels.SelectSession(db, sessionToken)
	if selectError != nil {
		if selectError.Error() == "sql: no rows in result set" {
			deleteCookie(w, "session_token")
			return false, userManagementModels.User{}, "", nil
		} else {
			return false, userManagementModels.User{}, "", selectError
		}
	}

	// Check if the cookie has expired
	if time.Now().After(expirationTime) {
		// Cookie expired, redirect to login
		return false, userManagementModels.User{}, "", nil
	}

	return true, user, sessionToken, nil
}

func Logout(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, loggedInUser, sessionToken, checkLoginError := CheckLogin(w, r, db)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	if !loginStatus {
		res := utils.Result{
			Success: true,
			Message: "You are not logged in",
			Data:    nil,
		}
		utils.ReturnJson(w, res)
		return
	}

	err := userManagementModels.DeleteSession(db, sessionToken)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	deleteCookie(w, "session_token") // Deleting a cookie named "session_token"
	Mutex.Lock()
	defer Mutex.Unlock()
	SocketLogoutHandler(w, r, loggedInUser)

	res := utils.Result{
		Success: true,
		Message: "Logged out successfully",
		Data:    nil,
	}
	utils.ReturnJson(w, res)
	return
}

func UpdateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, loginUser, _, checkLoginError := CheckLogin(w, r, db)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if !loginStatus {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
		return
	}

	const maxUploadSize = 2 << 20 // 2 MB

	// Limit the request body size
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	first_name := utils.SanitizeInput(r.FormValue("first_name"))
	last_name := utils.SanitizeInput(r.FormValue("last_name"))
	birthday_str := utils.SanitizeInput(r.FormValue("birthday"))
	gender := utils.SanitizeInput(r.FormValue("gender"))
	about_me := utils.SanitizeInput(r.FormValue("about_me"))
	visibility := utils.SanitizeInput(r.FormValue("visibility"))

	if len(first_name) == 0 || len(last_name) == 0 || len(birthday_str) == 0 || len(gender) == 0 {
		res := utils.Result{
			Success:    false,
			Message:    "first_name, last_name, birthday, gender are required.",
			HttpStatus: http.StatusOK,
			Data:       nil,
		}
		utils.ReturnJson(w, res)
		return
	}
	if len(visibility) == 0 {
		visibility = "private"
	}

	birthday, err := time.Parse("2006-01-02", birthday_str) // Adjust format as needed
	if err != nil {
		res := utils.Result{
			Success:    false,
			Message:    "Invalid birth date format. Use YYYY-MM-DD.",
			HttpStatus: http.StatusOK,
			Data:       nil,
		}
		utils.ReturnJson(w, res)
		return
	}

	profile_image_file, handler, err := r.FormFile("profile_image")
	if err != nil {
		// "File is missing"
		user := &userManagementModels.User{
			ID:         loginUser.ID,
			FirstName:  first_name,
			LastName:   last_name,
			Gender:     gender,
			BirthDay:   birthday,
			AboutMe:    about_me,
			Visibility: visibility,
		}

		// Update a record while checking duplicates
		updateError := userManagementModels.UpdateUser(db, user)
		if updateError != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		res := utils.Result{
			Success: true,
			Message: "Profile updated successfully",
			Data:    nil,
		}
		utils.ReturnJson(w, res)
		return
	} else {
		defer profile_image_file.Close()

		profile_image := ""
		if handler.Size != 0 {
			// Extra safety: check file size from the header
			if handler.Size > maxUploadSize {
				// "File is too large or missing"
				errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
				return
			}

			// Call your file upload function
			profile_image, err = fileManagementControllers.FileUpload(profile_image_file, handler)
			if err != nil {
				errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
				return
			}
		}

		user := &userManagementModels.User{
			ID:           loginUser.ID,
			FirstName:    first_name,
			LastName:     last_name,
			Gender:       gender,
			BirthDay:     birthday,
			AboutMe:      about_me,
			Visibility:   visibility,
			ProfileImage: profile_image,
		}

		// Update a record while checking duplicates
		updateError := userManagementModels.UpdateUser(db, user)
		if updateError != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		res := utils.Result{
			Success: true,
			Message: "Profile updated successfully",
			Data:    nil,
		}
		utils.ReturnJson(w, res)
		return
	}

}

func deleteCookie(w http.ResponseWriter, cookieName string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",              // Optional but recommended
		Expires:  time.Unix(0, 0), // Set expiration to a past date
		MaxAge:   -1,              // Ensure immediate removal
		Path:     "/",             // Must match the original cookie path
		HttpOnly: true,
		Secure:   false,
	})
}

func UserSetCookie(w http.ResponseWriter, sessionToken string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})
}

// Helper function to broadcast the list of online users
func UpdateOnlineUsers() {
	usernames := make([]string, 0, len(OnlineUsers))
	for _, username := range OnlineUsers {
		usernames = append(usernames, username)
	}

	// Encode the list of usernames as JSON
	userListJSON, err := json.Marshal(usernames)
	if err != nil {
		fmt.Println("Error encoding online users:", err)
		return
	}

	// Send the list to all online clients
	for client := range OnlineUsers {
		err := client.WriteMessage(websocket.TextMessage, userListJSON)
		if err != nil {
			Mutex.Lock()
			defer Mutex.Unlock()
			client.Close()
			delete(OnlineUsers, client)

		}
	}
}

func SocketLogoutHandler(w http.ResponseWriter, r *http.Request, loggedInUser userManagementModels.User) {
	for clientConn, clientUserName := range OnlineUsers {
		if clientUserName == loggedInUser.NickName {
			defer clientConn.Close() // close their websocket
			delete(OnlineUsers, clientConn)
			UpdateOnlineUsers() // Update the online users list
			fmt.Println("User logged out:", clientUserName, "| OnlineUsers:", OnlineUsers)
		}
		clientConn.WriteJSON(map[string]string{
			"type":    "logout",
			"message": "You have been logged out.",
		})
	}
}
