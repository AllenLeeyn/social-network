package controller

import (
	// "forum/middlewares"

	"net/http"
	errorControllers "social-network/pkg/errorManagement/controllers"
	fileControllers "social-network/pkg/fileManagement/controllers"
	"social-network/pkg/forumManagement/models"
	"social-network/pkg/utils"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func ReadAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("userID")
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

	utils.ReturnJsonSuccess(w, "Post fetched successfully", posts)
}

func ReadPostsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := models.ReadAllCategories()
	if err != nil {
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

	posts, err := models.ReadPostsByCategoryId(filteredCategory.ID)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		Posts                []models.Post
		Categories           []models.Category
		SelectedCategoryName string
	}{
		Posts:                posts,
		Categories:           categories,
		SelectedCategoryName: categoryName,
	}

	utils.ReturnJsonSuccess(w, "Post fetched successfully", data_obj_sender)
}

func FilterPostsHandler(w http.ResponseWriter, r *http.Request) {
	searchTerm, errUrl := utils.ExtractFromUrl(r.URL.Path, "api/filterPosts")
	if errUrl == "not found" {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}

	categories, err := models.ReadAllCategories()
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	posts, err := models.FilterPosts(searchTerm)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		Posts      []models.Post
		Categories []models.Category
		SearchTerm string
	}{
		Posts:      posts,
		Categories: categories,
		SearchTerm: searchTerm,
	}

	utils.ReturnJsonSuccess(w, "Post fetched successfully", data_obj_sender)
}

func ReadMyCreatedPostsHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("userID")
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	categories, err := models.ReadAllCategories()
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	posts, err := models.ReadPostsByUserId(userID)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		Posts      []models.Post
		Categories []models.Category
	}{
		Posts:      posts,
		Categories: categories,
	}

	utils.ReturnJsonSuccess(w, "Post fetched successfully", data_obj_sender)
}

func ReadMyLikedPostsHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("userID")
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	categories, err := models.ReadAllCategories()
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	posts, err := models.ReadPostsLikedByUserId(userID)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		Posts      []models.Post
		Categories []models.Category
	}{
		Posts:      posts,
		Categories: categories,
	}

	utils.ReturnJsonSuccess(w, "Post fetched successfully", data_obj_sender)
}

func ReadPostHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("userID")
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	uuid, errUrl := utils.ExtractUUIDFromUrl(r.URL.Path, "api/post")
	if errUrl == "not found" {
		errorControllers.ErrorHandler(w, r, errorControllers.NotFoundError)
		return
	}

	post, err := models.ReadPostByUUID(uuid, userID)
	if err != nil {
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

func SubmitPostHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("userID")
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	// Parse the multipart form with a max memory of 10MB
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}

	title := utils.SanitizeInput(r.FormValue("title"))
	content := utils.SanitizeInput(r.FormValue("content"))
	categories := r.Form["categories"]
	if len(title) == 0 || len(content) == 0 || len(categories) == 0 {
		res := utils.Result{
			Success:    false,
			Message:    "title, content and categories are required.",
			HttpStatus: http.StatusOK,
			Data:       nil,
		}
		utils.ReturnJson(w, res)
		return
	}

	// Retrieve all uploaded files
	files := r.MultipartForm.File["postFiles"]

	uploadedFiles := make(map[string]string)

	for _, handler := range files {
		file, err := handler.Open()
		if err != nil {
			errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
			return
		}
		defer file.Close()

		// Call your file upload function
		uploadedFile, err := fileControllers.FileUpload(file, handler)
		if err != nil {
			errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
			return
		}

		uploadedFiles[handler.Filename] = uploadedFile
	}

	post := &models.Post{
		Title:   title,
		Content: content,
		UserId:  userID,
	}

	// Convert the string slice to an int slice
	categoryIds := make([]int, 0, len(categories))
	for _, category := range categories {
		if id, err := strconv.Atoi(category); err == nil {
			categoryIds = append(categoryIds, id)
		} else {
			// Handle error if conversion fails (for example, invalid input)
			errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
			return
		}
	}

	// Insert a record while checking duplicates
	_, insertError := models.InsertPost(post, categoryIds, uploadedFiles)
	if insertError != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Post submitted successfully", nil)
}

func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("userID")
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	// Parse the multipart form with a max memory of 10MB
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}

	idStr := r.FormValue("post_id")
	// uuid := utils.SanitizeInput(r.FormValue("uuid"))
	title := utils.SanitizeInput(r.FormValue("title"))
	content := utils.SanitizeInput(r.FormValue("content"))
	categories := r.Form["update_post_categories"]

	if len(idStr) == 0 || len(title) == 0 || len(content) == 0 || len(categories) == 0 {
		res := utils.Result{
			Success:    false,
			Message:    "title, description and categories are required.",
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

	// Retrieve all uploaded files
	files := r.MultipartForm.File["postFiles"]

	uploadedFiles := make(map[string]string)

	for _, handler := range files {
		file, err := handler.Open()
		if err != nil {
			errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
			return
		}
		defer file.Close()

		// Call your file upload function
		uploadedFile, err := fileControllers.FileUpload(file, handler)
		if err != nil {
			errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
			return
		}

		uploadedFiles[handler.Filename] = uploadedFile
	}

	post := &models.Post{
		ID:      id,
		Title:   title,
		Content: content,
		UserId:  userID,
	}

	// Convert the string slice to an int slice
	categoryIds := make([]int, 0, len(categories))
	for _, category := range categories {
		if id, err := strconv.Atoi(category); err == nil {
			categoryIds = append(categoryIds, id)
		} else {
			// Handle error if conversion fails (for example, invalid input)
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
	}

	// Update a record while checking duplicates
	updateError := models.UpdatePost(post, categoryIds, uploadedFiles, userID)
	if updateError != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Post updated successfully", nil)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
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

	idStr := r.FormValue("id")

	if len(idStr) == 0 {
		errorControllers.ErrorHandler(w, r, errorControllers.BadRequestError)
		return
	}

	post_id, err := strconv.Atoi(idStr)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	// Update a record while checking duplicates
	updateError := models.UpdateStatusPost(post_id, "delete", userID)
	if updateError != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Post removed successfully", nil)
}

func PostFeedbackHandler(w http.ResponseWriter, r *http.Request) {
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
	postID := r.FormValue("post_id")
	postIDInt, _ := strconv.Atoi(postID)
	// var Type string
	// like := r.FormValue("like_post")
	// dislike := r.FormValue("dislike_post")
	// if like == "like" {
	// 	Type = like
	// } else if dislike == "dislike" {
	// 	Type = dislike
	// }
	ratingStr := r.FormValue("rating")

	existingLikeId, existingFeedbackRating := models.PostHasFeedbacked(userID, postIDInt)

	var resMessage string
	if ratingStr == "1" {
		resMessage = "You liked successfully"
	} else {
		resMessage = "You disliked successfully"
	}

	rating, err := strconv.Atoi(ratingStr)
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	if existingLikeId == -1 {
		post := &models.PostFeedback{
			Rating: rating,
			PostId: postIDInt,
			UserId: userID,
		}
		_, insertError := models.InsertPostFeedback(post)
		if insertError != nil {
			errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
			return
		}

		utils.ReturnJsonSuccess(w, resMessage, nil)
		return
	} else {
		updateError := models.UpdateStatusFeedback(existingLikeId, "delete", userID)
		if updateError != nil {
			errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
			return
		}

		if existingFeedbackRating != ratingStr { //this is duplicated like or duplicated dislike so we should update it to disable
			post := &models.PostFeedback{
				Rating: rating,
				PostId: postIDInt,
				UserId: userID,
			}
			_, insertError := models.InsertPostFeedback(post)
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
