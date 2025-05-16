package controller

import (
	// "forum/middlewares"

	"errors"
	"net/http"
	errorControllers "social-network/pkg/errorManagement/controllers"
	"social-network/pkg/forumManagement/models"
	"social-network/pkg/middleware"
	"social-network/pkg/utils"

	_ "github.com/mattn/go-sqlite3"
)

func isValidPostFeedback(postFeedback *models.PostFeedback) error {
	isValid := false

	if postFeedback.ParentId, isValid = utils.IsValidId(postFeedback.ParentId); !isValid {
		return errors.New("parent id is required and must be numeric")
	}
	if postFeedback.Rating, isValid = utils.IsValidRating(postFeedback.Rating); !isValid {
		return errors.New("rating is required and must be 1 or -1")
	}
	return nil
}

func PostFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	postFeedback := &models.PostFeedback{}
	postFeedback.UserId = userID
	if err := utils.ReadJSON(w, r, postFeedback); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	if err := isValidPostFeedback(postFeedback); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	existingFeedbackRating := models.PostHasFeedback(userID, postFeedback.ParentId)

	var resMessage string
	if postFeedback.Rating == 1 {
		resMessage = "You liked successfully"
	} else {
		resMessage = "You disliked successfully"
	}

	if existingFeedbackRating == -1000 {
		_, insertError := models.InsertPostFeedback(postFeedback)
		if insertError != nil {
			errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
			return
		}

		utils.ReturnJsonSuccess(w, resMessage, nil)
		return
	} else {
		if existingFeedbackRating != postFeedback.Rating {
			updateError := models.UpdatePostFeedback(postFeedback)
			if updateError != nil {
				errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
				return
			}
		} else { //this is duplicated like or duplicated dislike so we should update it to delete
			if postFeedback.Rating == 1 {
				resMessage = "You removed like successfully"
			} else {
				resMessage = "You removed dislike successfully"
			}

			postFeedback.Rating = 0
			updateError := models.UpdatePostFeedback(postFeedback)
			if updateError != nil {
				errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
				return
			}
		}
		utils.ReturnJsonSuccess(w, resMessage, nil)
		return
	}
}
