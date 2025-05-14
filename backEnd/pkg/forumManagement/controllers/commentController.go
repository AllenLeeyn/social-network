package controller

import (
	"net/http"
	errorControllers "social-network/pkg/errorManagement/controllers"
	"social-network/pkg/forumManagement/models"
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
	userIDRaw := r.Context().Value("userID")
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

func SubmitCommentHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("userID")
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	err := r.ParseMultipartForm(0)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}

	post_id_str := r.FormValue("post_id")
	content := utils.SanitizeInput(r.FormValue("content"))
	if len(post_id_str) == 0 || len(content) == 0 {
		res := utils.Result{
			Success:    false,
			Message:    "content is required.",
			HttpStatus: http.StatusOK,
			Data:       nil,
		}
		utils.ReturnJson(w, res)
		return
	}

	post_id, err := strconv.Atoi(post_id_str)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	// Insert a record while checking duplicates
	_, insertError := models.InsertComment(post_id, userID, content)
	if insertError != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Comment submitted successfully", nil)
}

func FeedbackCommentHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("userID")
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

func UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("userID")
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	err := r.ParseMultipartForm(0)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}

	idStr := r.FormValue("comment_id")
	post_uuid := utils.SanitizeInput(r.FormValue("post_uuid"))
	content := utils.SanitizeInput(r.FormValue("content"))

	if len(idStr) == 0 || len(post_uuid) == 0 || len(content) == 0 {
		res := utils.Result{
			Success:    false,
			Message:    "content is required.",
			HttpStatus: http.StatusOK,
			Data:       nil,
		}
		utils.ReturnJson(w, res)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	comment := &models.Comment{
		ID:      id,
		Content: content,
		UserId:  userID,
	}

	// Update a record while checking duplicates
	updateError := models.UpdateComment(comment, userID, content)
	if updateError != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Comment updated successfully", nil)
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("userID")
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	err := r.ParseMultipartForm(0)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}

	idStr := r.FormValue("comment_id")
	post_uuid := utils.SanitizeInput(r.FormValue("post_uuid"))

	if len(idStr) == 0 || len(post_uuid) == 0 {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}

	comment_id, err := strconv.Atoi(idStr)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	// Update a record while checking duplicates
	updateError := models.UpdateCommentStatus(comment_id, "delete", userID)
	if updateError != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Comment removed successfully", nil)
}
