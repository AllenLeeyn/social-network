package handlers

import (
	"net/http"
	"strconv"
)

// Post() to grab individual post data.
func Post(w http.ResponseWriter, r *http.Request) {
	sessionCookie, userID, isValid := checkHttpRequest(w, r, "user", http.MethodGet)
	if !isValid {
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		executeJSON(w, MsgData{"Invalid post ID"}, http.StatusBadRequest)
		return
	}
	selectedPost, err := db.SelectPost(id, userID)
	if err != nil {
		executeJSON(w, MsgData{"Error getting post"}, http.StatusBadRequest)
		return
	}
	if selectedPost == nil {
		executeJSON(w, MsgData{"Post not found"}, http.StatusNotFound)
		return
	}
	seletctedComments, err := db.SelectComments(id, userID, "")
	if err != nil {
		executeJSON(w, MsgData{"Error getting comments"}, http.StatusInternalServerError)
		return
	}
	extendSession(w, sessionCookie)
	executeJSON(w, postData{
		Post:     *selectedPost,
		Comments: *seletctedComments,
	}, http.StatusOK)
}
