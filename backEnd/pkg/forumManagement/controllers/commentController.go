package controller

import (
	"errors"
	"fmt"
	"net/http"
	errorControllers "social-network/pkg/errorManagement/controllers"
	"social-network/pkg/forumManagement/models"
	"social-network/pkg/middleware"
	"social-network/pkg/utils"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func ReadAllCommentsHandler(w http.ResponseWriter, r *http.Request) {
	comments, err := models.ReadAllComments()
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Comments fetched successfully", comments)
}

func ReadPostCommentsHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	// todo get post id from url start
	postIdStr := r.URL.Query().Get("post_id")
	if len(postIdStr) == 0 {
		res := utils.Result{
			Success:    false,
			Message:    "post_id is required.",
			HttpStatus: http.StatusOK,
			Data:       nil,
		}
		utils.ReturnJson(w, res)
		return
	}
	//todo get post id from url end
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}
	comments, err := models.ReadAllCommentsOfUserForPost(postId, userID)
	// comments, err := models.ReadAllCommentsForPost(userID)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Comments fetched successfully", comments)
}

func isValidCommentInfo(comment *models.Comment) error {
	isValid := false

	if comment.PostId, isValid = utils.IsValidId(comment.PostId); !isValid {
		return errors.New("post id is required and must be numeric")
	}
	if comment.Content, isValid = utils.IsValidContent(comment.Content, 3, 1000); !isValid {
		return errors.New("content is required and must be between 3 to 1000 alphanumeric characters, '_' or '-'")
	}

	//todo check parent_id

	return nil
}

// todo
func SubmitCommentHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	comment := &models.Comment{}
	comment.UserId = userID
	if err := utils.ReadJSON(w, r, comment); err != nil {
		fmt.Println("Error reading JSON:", err)
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	if err := isValidCommentInfo(comment); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert a record while checking duplicates
	_, insertError := models.InsertComment(comment)
	if insertError != nil {
		fmt.Println("Error inserting comment:", insertError)
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Comment submitted successfully", nil)
}

func isValidUpdateCommentInfo(comment *models.Comment) error {
	isValid := false

	if comment.ID, isValid = utils.IsValidId(comment.ID); !isValid {
		return errors.New("comment id is required and must be numeric")
	}

	return isValidCommentInfo(comment)
}

// todo
func UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	comment := &models.Comment{}
	comment.UserId = userID
	if err := utils.ReadJSON(w, r, comment); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	if err := isValidUpdateCommentInfo(comment); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	// Update a record while checking duplicates
	updateError := models.UpdateComment(comment)
	if updateError != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Comment updated successfully", nil)
}

// todo
func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	comment := &models.Comment{}
	if err := utils.ReadJSON(w, r, comment); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	if _, isValid := utils.IsValidId(comment.ID); !isValid {
		errorControllers.CustomErrorHandler(w, r, "comment id is required and must be numeric", http.StatusBadRequest)
		return
	}

	// Update a record while checking duplicates
	updateError := models.UpdateCommentStatus(comment.ID, "delete", userID)
	if updateError != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Comment removed successfully", nil)
}
