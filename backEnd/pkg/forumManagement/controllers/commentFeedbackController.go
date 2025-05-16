package controller

import (
	"net/http"
	errorControllers "social-network/pkg/errorManagement/controllers"
	"social-network/pkg/forumManagement/models"
	"social-network/pkg/middleware"
	"social-network/pkg/utils"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func CommentFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	// err := r.ParseForm()
	err := r.ParseMultipartForm(0)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}
	parentID := r.FormValue("parent_id")
	parentIDInt, _ := strconv.Atoi(parentID)
	ratingStr := r.FormValue("rating")

	rating, err := strconv.Atoi(ratingStr)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}

	existingFeedbackId, existingFeedbackRating := models.CommentHasFeedback(userID, parentIDInt)

	var resMessage string
	if ratingStr == "1" {
		resMessage = "You liked successfully"
	} else {
		resMessage = "You disliked successfully"
	}

	if existingFeedbackId == -1 {
		insertError := models.InsertCommentFeedback(ratingStr, parentIDInt, userID)
		if insertError != nil {
			errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
			return
		}

		utils.ReturnJsonSuccess(w, resMessage, nil)
		return
	} else {
		updateError := models.UpdateCommentFeedbackStatus(existingFeedbackId, "delete", userID)
		if updateError != nil {
			errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
			return
		}

		if existingFeedbackRating != rating { //this is duplicated like or duplicated dislike so we should update it to disable
			insertError := models.InsertCommentFeedback(ratingStr, parentIDInt, userID)
			if insertError != nil {
				errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
				return
			}
		} else {
			if ratingStr == "1" {
				resMessage = "You removed like successfully"
			} else {
				resMessage = "You removed dislike successfully"
			}
		}

		utils.ReturnJsonSuccess(w, resMessage, nil)
		return
	}
}
