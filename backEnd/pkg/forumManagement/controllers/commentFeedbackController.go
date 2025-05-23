package controller

import (
	"errors"
	"net/http"
	errorControllers "social-network/pkg/errorManagement/controllers"
	"social-network/pkg/forumManagement/models"
	"social-network/pkg/middleware"
	"social-network/pkg/utils"

	_ "github.com/mattn/go-sqlite3"
)

func isValidCommentFeedback(commentFeedback *models.CommentFeedback) error {
	isValid := false

	if commentFeedback.ParentId, isValid = utils.IsValidId(commentFeedback.ParentId); !isValid {
		return errors.New("parent id is required and must be numeric")
	}
	if commentFeedback.Rating, isValid = utils.IsValidRating(commentFeedback.Rating); !isValid {
		return errors.New("rating is required and must be 1 or -1")
	}
	return nil
}

func CommentFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	commentFeedback := &models.CommentFeedback{}
	commentFeedback.UserId = userID
	if err := utils.ReadJSON(r, commentFeedback); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	if err := isValidCommentFeedback(commentFeedback); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	existingFeedbackRating := models.CommentHasFeedback(userID, commentFeedback.ParentId)

	var resMessage string
	if commentFeedback.Rating == 1 {
		resMessage = "You liked successfully"
	} else {
		resMessage = "You disliked successfully"
	}

	if existingFeedbackRating == -1000 {
		lastInsertID, insertError := models.InsertCommentFeedback(commentFeedback)
		if insertError != nil {
			errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
			return
		}

		utils.ReturnJsonSuccess(w, resMessage, lastInsertID)
		return
	} else {
		if existingFeedbackRating != commentFeedback.Rating {
			updateError := models.UpdateCommentFeedback(commentFeedback)
			if updateError != nil {
				errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
				return
			}
		} else { //this is duplicated like or duplicated dislike so we should update it to delete
			if commentFeedback.Rating == 1 {
				resMessage = "You removed like successfully"
			} else {
				resMessage = "You removed dislike successfully"
			}

			commentFeedback.Rating = 0
			updateError := models.UpdateCommentFeedback(commentFeedback)
			if updateError != nil {
				errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
				return
			}
		}
		utils.ReturnJsonSuccess(w, resMessage, nil)
		return
	}
}
