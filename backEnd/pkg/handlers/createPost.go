package handlers

import (
	"net/http"
	"strconv"
)

// CreatePost() and redirects to new post if session is valid
func CreatePost(w http.ResponseWriter, r *http.Request) {
	sessionCookie, userID, isValid := checkHttpRequest(w, r, "user", http.MethodPost)
	if !isValid {
		return
	}

	data := &post{}
	if !getJSON(w, r, data) {
		return
	}

	isValidTitle, title := checkContent(data.Title, 10, 200)
	isValidContent, content := checkContent(data.Content, 10, 2000)
	if !isValidTitle || !isValidContent {
		executeJSON(w, MsgData{"Error: " + title + content}, http.StatusBadRequest)
		return
	}
	data.UserID = userID
	data.Title = title
	data.Content = content
	postNum, err := db.InsertPost(data)
	if err != nil {
		executeJSON(w,
			MsgData{"Error creating post. Must select at least one category."},
			http.StatusInternalServerError)
		return
	}
	postID := strconv.Itoa(postNum)
	extendSession(w, sessionCookie)
	executeJSON(w, MsgData{postID}, http.StatusOK)

}
