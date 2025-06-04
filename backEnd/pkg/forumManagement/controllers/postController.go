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

// validators for post id
func isValidPostId(post *models.Post) error {
	isValid := false

	if post.ID, isValid = utils.IsValidId(post.ID); !isValid {
		return errors.New("comment id is required and must be numeric")
	}

	return nil
}

// validators for post data
func isValidPostInfo(post *models.Post) error {
	isValid := false

	if post.Title, isValid = utils.IsValidContent(post.Title, 3, 100); !isValid {
		return errors.New("title is required and must be between 3 to 100 alphanumeric characters, '_' or '-'")
	}
	if post.Content, isValid = utils.IsValidContent(post.Content, 3, 1000); !isValid {
		return errors.New("content is required and must be between 3 to 1000 alphanumeric characters, '_' or '-'")
	}
	// if post.CategoryIds, isValid = utils.IsValidCategoryIdsList(post.CategoryIds); !isValid {
	if post.CategoryIds, isValid = utils.IsValidIntegerList(post.CategoryIds); !isValid {
		return errors.New("category is required and must be a list of integers")
	}

	if post.Visibility == "selected" {
		if post.SelectedAudienceUserUUIDS == nil {
			return errors.New("no audience is selected")
		}
	}
	if post.Visibility == "" {
		post.Visibility = "public"
	}
	if post.Visibility == "public" {
		post.GroupId = 0
	}
	if post.Visibility != "public" && post.Visibility != "private" && post.Visibility != "selected" {
		return errors.New("visibility is required and must be either 'public', 'private' or 'selected'")
	}

	if post.Type == "" {
		post.Type = "user"
	} else if post.Type != "group" && post.Type != "user" {
		return errors.New("type is required and must be either 'group' or 'user'")
	}

	if post.Type == "user" && post.GroupId != 0 {
		return errors.New("group id is not required for user type")
	} else if post.Type == "group" && post.GroupId == 0 {
		return errors.New("group id is required and must be numeric")
	} else if post.Type == "group" && post.GroupId != 0 {
		if post.GroupId, isValid = utils.IsValidId(post.GroupId); !isValid {
			return errors.New("group id is required and must be numeric")
		}
	}
	return nil
}

// validators for update post data
func isValidUpdatePostInfo(post *models.Post) error {
	isValid := false

	if post.ID, isValid = utils.IsValidId(post.ID); !isValid {
		return errors.New("post id is required and must be numeric")
	}

	if errValidPostId := isValidPostId(post); errValidPostId != nil {
		return errValidPostId
	}

	if errisValidPostInfo := isValidPostInfo(post); errisValidPostInfo != nil {
		return errisValidPostInfo
	}
	return nil
}

func ReadAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	posts, err := models.ReadAllPosts(userID)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Posts fetched successfully", posts)
}

func ReadPostsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	categoryName, errUrl := utils.ExtractFromUrl(r.URL.Path, "api/posts")
	if errUrl == "not found" {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}

	filteredCategory, errCategory := models.ReadCategoryByName(categoryName)
	if errCategory != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}

	posts, err := models.ReadPostsByCategoryId(filteredCategory.ID, userID)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		Posts                []models.Post
		SelectedCategoryName string
		SelectedCategoryId   int
	}{
		Posts:                posts,
		SelectedCategoryName: categoryName,
		SelectedCategoryId:   filteredCategory.ID,
	}

	utils.ReturnJsonSuccess(w, "Posts fetched successfully", data_obj_sender)
}

func FilterPostsHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	searchTerm, errUrl := utils.ExtractFromUrl(r.URL.Path, "api/filterPosts")
	if errUrl == "not found" {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}

	posts, err := models.FilterPosts(searchTerm, userID)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		Posts      []models.Post
		SearchTerm string
	}{
		Posts:      posts,
		SearchTerm: searchTerm,
	}

	utils.ReturnJsonSuccess(w, "Posts fetched successfully", data_obj_sender)
}

func ReadMyCreatedPostsHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	posts, err := models.ReadPostsByUserId(userID)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Posts fetched successfully", posts)
}

func ReadMyLikedPostsHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	posts, err := models.ReadPostsLikedByUserId(userID)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Posts fetched successfully", posts)
}

func ReadPostHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	uuid, err := utils.ExtractUUIDFromUrl(r.URL.Path, "api/post")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}

	post, err := models.ReadPostByUUID(uuid, userID)
	if err != nil {
		if err.Error() == "uuid not found" {
			errorControllers.CustomErrorHandler(w, r, "post with this uuid not found for you", errorControllers.NotFoundError.CodeNumber)
			return
		}
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		Post     models.Post
		Comments []models.Comment
	}{
		Post:     post,
		Comments: nil,
	}

	comments, err := models.ReadAllCommentsForPostByUserID(post.ID, userID)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	data_obj_sender.Comments = comments

	utils.ReturnJsonSuccess(w, "Post fetched successfully", data_obj_sender)
}

func ReadPostsSubmittedByUserUUIDHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	audienceUserId, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	userUUID, err := utils.ExtractUUIDFromUrl(r.URL.Path, "api/userPosts")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}

	posts, err := models.ReadPostsSubmittedByUserUUID(userUUID, audienceUserId)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Posts fetched successfully", posts)
}

func ReadPostsSubmittedByGroupUUIDHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	audienceUserId, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	groupUUID, err := utils.ExtractUUIDFromUrl(r.URL.Path, "api/groupPosts")
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}

	posts, err := models.ReadPostsSubmittedByGroupUUID(groupUUID, audienceUserId)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Posts fetched successfully", posts)
}

func SubmitPostHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	post := &models.Post{}
	post.UserId = userID
	if err := utils.ReadJSON(r, post); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	if err := isValidPostInfo(post); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert a record while checking duplicates
	createdPostUUID, insertError := models.InsertPost(post, post.CategoryIds, post.FileAttachments)
	if insertError != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Post submitted successfully", createdPostUUID)
}

func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	post := &models.Post{}
	post.UserId = userID
	if err := utils.ReadJSON(r, post); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	if err := isValidUpdatePostInfo(post); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	// Update a record while checking duplicates
	updateError := models.UpdatePost(post, post.CategoryIds, post.FileAttachments, userID)
	if updateError != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Post updated successfully", nil)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middleware.CtxUserID)
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	post := &models.Post{}
	if err := utils.ReadJSON(r, post); err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	if err := isValidPostId(post); err != nil {
		errorControllers.CustomErrorHandler(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	updateError := models.UpdateStatusPost(post.ID, "delete", userID)
	if updateError != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Post removed successfully", nil)
}
