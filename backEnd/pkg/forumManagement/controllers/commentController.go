package controller

import (
	"database/sql"
	"net/http"
	errorManagementControllers "social-network/pkg/errorManagement/controllers"
	"social-network/pkg/forumManagement/models"
	"social-network/pkg/utils"
	"strconv"

	userManagementControllers "social-network/pkg/userManagement/controllers"

	_ "github.com/mattn/go-sqlite3"
)

func ReadAllComments(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, _, _, checkLoginError := userManagementControllers.CheckLogin(w, r, db)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if !loginStatus {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
		return
	}

	comments, err := models.ReadAllComments(db)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	res := utils.Result{
		Success: true,
		Message: "Comments fetched successfully",
		Data:    comments,
	}
	utils.ReturnJson(w, res)
}

func ReadPostComments(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, loginUser, _, checkLoginError := userManagementControllers.CheckLogin(w, r, db)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if !loginStatus {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
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
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	comments, err := models.ReadAllCommentsOfUserForPost(db, postId, loginUser.ID)
	// comments, err := models.ReadAllCommentsForPost(db, loginUser.ID)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	res := utils.Result{
		Success: true,
		Message: "Comments fetched successfully",
		Data:    comments,
	}
	utils.ReturnJson(w, res)
}

func SubmitComment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, loginUser, _, checkLoginError := userManagementControllers.CheckLogin(w, r, db)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if !loginStatus {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
		return
	}

	err := r.ParseMultipartForm(0)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
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
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	// Insert a record while checking duplicates
	_, insertError := models.InsertComment(db, post_id, loginUser.ID, content)
	if insertError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	res := utils.Result{
		Success: true,
		Message: "Comment submitted successfully",
		Data:    nil,
	}
	utils.ReturnJson(w, res)
}

func FeedbackComment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, loginUser, _, checkLoginError := userManagementControllers.CheckLogin(w, r, db)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if !loginStatus {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
		return
	}

	// err := r.ParseForm()
	err := r.ParseMultipartForm(0)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}
	parentID := r.FormValue("parent_id")
	parentIDInt, _ := strconv.Atoi(parentID)
	ratingStr := r.FormValue("rating")

	rating, err := strconv.Atoi(ratingStr)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	existingFeedbackId, existingFeedbackRating := models.CommentHasFeedback(db, loginUser.ID, parentIDInt)

	var resMessage string
	if ratingStr == "1" {
		resMessage = "You liked successfully"
	} else {
		resMessage = "You disliked successfully"
	}

	if existingFeedbackId == -1 {
		insertError := models.InsertCommentFeedback(db, ratingStr, parentIDInt, loginUser.ID)
		if insertError != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		res := utils.Result{
			Success: true,
			Message: resMessage,
			Data:    nil,
		}
		utils.ReturnJson(w, res)
		return
	} else {
		updateError := models.UpdateCommentFeedbackStatus(db, existingFeedbackId, "delete", loginUser.ID)
		if updateError != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		if existingFeedbackRating != rating { //this is duplicated like or duplicated dislike so we should update it to disable
			insertError := models.InsertCommentFeedback(db, ratingStr, parentIDInt, loginUser.ID)
			if insertError != nil {
				errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
				return
			}
		} else {
			if ratingStr == "1" {
				resMessage = "You removed like successfully"
			} else {
				resMessage = "You removed dislike successfully"
			}
		}
		res := utils.Result{
			Success: true,
			Message: resMessage,
			Data:    nil,
		}
		utils.ReturnJson(w, res)
		return
	}
}

func UpdateComment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, loginUser, _, checkLoginError := userManagementControllers.CheckLogin(w, r, db)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if !loginStatus {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
		return
	}

	err := r.ParseMultipartForm(0)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
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
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	comment := &models.Comment{
		ID:      id,
		Content: content,
		UserId:  loginUser.ID,
	}

	// Update a record while checking duplicates
	updateError := models.UpdateComment(db, comment, loginUser.ID, content)
	if updateError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	res := utils.Result{
		Success: true,
		Message: "Comment updated successfully",
		Data:    nil,
	}
	utils.ReturnJson(w, res)
}

func DeleteComment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, loginUser, _, checkLoginError := userManagementControllers.CheckLogin(w, r, db)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	if !loginStatus {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
		return
	}

	err := r.ParseMultipartForm(0)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	idStr := r.FormValue("comment_id")
	post_uuid := utils.SanitizeInput(r.FormValue("post_uuid"))

	if len(idStr) == 0 || len(post_uuid) == 0 {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	comment_id, err := strconv.Atoi(idStr)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	// Update a record while checking duplicates
	updateError := models.UpdateCommentStatus(db, comment_id, "delete", loginUser.ID)
	if updateError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	res := utils.Result{
		Success: true,
		Message: "Comment removed successfully",
		Data:    nil,
	}
	utils.ReturnJson(w, res)
}
