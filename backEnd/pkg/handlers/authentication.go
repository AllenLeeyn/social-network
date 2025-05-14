package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	_, _, isValid := checkHttpRequest(w, r, "guest", http.MethodPost)
	if !isValid {
		return
	}

	u := &user{}
	if !getJSON(w, r, u) {
		return
	}

	fmt.Println("Signup request received for:", u)

	if u.Gender == "" {
		u.Gender = "Other"
	}

	// check that credentials are valid
	e := checkCredentials(u)
	if e != nil {
		executeJSON(w, MsgData{e.Error()}, http.StatusBadRequest)
		return
	}
	passwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Passwd), 0)
	if err != nil {
		executeJSON(w, MsgData{"Error creating user"}, http.StatusInternalServerError)
		return
	}

	u.TypeID = 1
	u.PwHash = passwdHash
	u.ID, err = db.InsertUser(u)
	if err != nil {
		executeJSON(w, MsgData{"Error creating user"}, http.StatusInternalServerError)
		return
	}
	createSession(w, u)
	executeJSON(w, MsgData{"Signup succesful"}, http.StatusOK)
}

func Login(w http.ResponseWriter, r *http.Request) {
	_, _, isValid := checkHttpRequest(w, r, "guest", http.MethodPost)
	if !isValid {
		return
	}

	u := &user{}
	if !getJSON(w, r, u) {
		return
	}
	fmt.Println("hello")
	fmt.Println("Login request received for:", u.Email, u.NickName, u.Passwd) // Debug: Log user credentials

	if err := checkLoginCredentials(u); err != nil {
		executeJSON(w, MsgData{err.Error()}, http.StatusBadRequest)
		return
	}

	// check that user exists
	user, _ := db.SelectUserByField("nick_name", u.NickName)
	if user == nil {
		user, _ = db.SelectUserByField("email", u.Email)
	}
	if user == nil || bcrypt.CompareHashAndPassword(user.PwHash, []byte(u.Passwd)) != nil {
		executeJSON(w, MsgData{"Incorrect username and/or password"}, http.StatusBadRequest)
		return
	}

	s, e := db.SelectActiveSessionBy("user_id", user.ID)
	if e == nil {
		expireSession(w, s)
	}
	createSession(w, user)
	executeJSON(w, MsgData{"Login succesful"}, http.StatusOK)
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	sessionCookie, _ := r.Cookie("session-id")
	if sessionCookie == nil {
		executeJSON(w, MsgData{"You're not logged in"}, http.StatusBadRequest)
		return
	} else {
		s, _ := db.SelectActiveSessionBy("id", sessionCookie.Value)
		expireSession(w, s)
	}
	http.Redirect(w, r, "./", http.StatusSeeOther)
}

func checkLoginCredentials(u *user) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	usernameRegex := `^[a-zA-Z0-9_-]{3,16}$`

	if !validRegex(u.NickName, usernameRegex) {
		u.NickName = ""
	}
	if !validRegex(u.Email, emailRegex) {
		u.Email = ""
	}
	if !validPsswrd(u.Passwd) {
		return errors.New("password must be 8 characters or longer.\n" +
			"Include at least a lower case character, an upper case character, a number and one of '@$!%*?&'")
	}
	return nil
}

/*----------- authenticaton nefore handshake-----------*/
func WS(w http.ResponseWriter, r *http.Request) {
	sessionCookie, userID, isValid := checkHttpRequest(w, r, "user", http.MethodGet)
	if !isValid {
		return
	}
	u, err := db.SelectUserByField("id", userID)
	if err != nil || u == nil {
		executeJSON(w, MsgData{"User not found."}, http.StatusBadRequest)
		return
	}
	m.WebSocketUpgrade(w, r, sessionCookie.Value, u)
}
