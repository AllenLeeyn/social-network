package handlers

import (
	"net/http"
)

func CreateFeedback(w http.ResponseWriter, r *http.Request) {
	sessionCookie, userID, isValid := checkHttpRequest(w, r, "user", http.MethodPost)
	if !isValid {
		return
	}

	data := &feedback{}
	if !getJSON(w, r, data) {
		return
	}

	fb, err := db.SelectFeedback(data.Tgt, userID, data.ParentID)
	if err != nil {
		executeJSON(w, MsgData{"Error reading feedback"}, http.StatusInternalServerError)
		return
	} else if fb == nil {
		data.UserID = userID
		fb = data
		if err = db.InsertFeedback(data.Tgt, fb); err != nil {
			println(err.Error())
			executeJSON(w, MsgData{"Error reading feedback"}, http.StatusInternalServerError)
			return
		}
	} else {
		fb.Rating = data.Rating
		if err = db.UpdateFeedback(data.Tgt, fb); err != nil {
			executeJSON(w, MsgData{"Error reading feedback"}, http.StatusInternalServerError)
			return
		}
	}
	extendSession(w, sessionCookie)
	w.WriteHeader(http.StatusOK)
}
