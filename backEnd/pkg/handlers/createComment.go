package handlers

import (
	"net/http"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	sessionCookie, userID, isValid := checkHttpRequest(w, r, "user", http.MethodPost)
	if !isValid {
		return
	}

	data := &comment{}
	if !getJSON(w, r, data) {
		return
	}

	isValidContent, content := checkContent(data.Content, 10, 2000)
	if !isValidContent {
		executeJSON(w, MsgData{content}, http.StatusInternalServerError)
		return
	}

	post, err := db.SelectPost(data.PostID, userID)
	if err != nil || post == nil {
		executeJSON(w, MsgData{"Post not found"}, http.StatusNotFound)
		return
	}

	user, err := db.SelectUserByField("id", userID)
	if err != nil || user == nil {
		executeJSON(w, MsgData{"User not found"}, http.StatusNotFound)
		return
	}
	data.UserID = userID
	data.UserName = user.NickName
	data.Content = content

	if err := db.InsertComment(data); err != nil {
		executeJSON(w, MsgData{"Error creating comment"}, http.StatusInternalServerError)
		return
	}
	executeJSON(w, MsgData{"Comment created"}, http.StatusOK)
	extendSession(w, sessionCookie)
}
