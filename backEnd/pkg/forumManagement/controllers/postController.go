package controller

import (
	// "forum/middlewares"
	"net/http"
	"social-network/pkg/dbTools"
	errorManagementControllers "social-network/pkg/errorManagement/controllers"
	fileManagementControllers "social-network/pkg/fileManagement/controllers"
	"social-network/pkg/forumManagement/models"
	"social-network/pkg/utils"
	"strconv"
	"text/template"

	userManagementControllers "social-network/pkg/userManagement/controllers"
	userManagementModels "social-network/pkg/userManagement/models"

	_ "github.com/mattn/go-sqlite3"
)

func ReadAllPosts(w http.ResponseWriter, r *http.Request, db *dbTools.DBContainer) {
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

	posts, err := models.ReadAllPosts(db, loginUser.ID)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	res := utils.Result{
		Success: true,
		Message: "Posts fetched successfully",
		Data:    posts,
	}
	utils.ReturnJson(w, res)
}

func AdminReadAllPosts(w http.ResponseWriter, r *http.Request, db *dbTools.DBContainer) {
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

	posts, err := models.ReadAllPosts(db, loginUser.ID)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		LoginUser userManagementModels.User
		Posts     []models.Post
	}{
		LoginUser: userManagementModels.User{},
		Posts:     posts,
	}

	if loginStatus {
		data_obj_sender.LoginUser = loginUser
	}

	// Create a template with a function map
	tmpl, err := template.New("admin_posts.html").Funcs(template.FuncMap{
		"formatDate": utils.FormatDate, // Register function globally
	}).ParseFiles(
		publicUrl+"admin_posts.html",
		publicUrl+"templates/footer.html",
	)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	err = tmpl.Execute(w, data_obj_sender)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
}

func ReadPostsByCategory(w http.ResponseWriter, r *http.Request, db *dbTools.DBContainer) {
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

	categories, err := models.ReadAllCategories(db)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	categoryName, errUrl := utils.ExtractFromUrl(r.URL.Path, "api/posts")
	if errUrl == "not found" {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.NotFoundError)
		return
	}

	filteredCategory, errCategory := models.ReadCategoryByName(db, categoryName)
	if errCategory != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.NotFoundError)
		return
	}

	posts, err := models.ReadPostsByCategoryId(db, filteredCategory.ID)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		LoginUser            userManagementModels.User
		Posts                []models.Post
		Categories           []models.Category
		SelectedCategoryName string
	}{
		LoginUser:            userManagementModels.User{},
		Posts:                posts,
		Categories:           categories,
		SelectedCategoryName: categoryName,
	}

	if loginStatus {
		data_obj_sender.LoginUser = loginUser
	}

	res := utils.Result{
		Success: true,
		Message: "Posts fetched successfully",
		Data:    data_obj_sender,
	}
	utils.ReturnJson(w, res)
}

func FilterPosts(w http.ResponseWriter, r *http.Request, db *dbTools.DBContainer) {
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

	searchTerm, errUrl := utils.ExtractFromUrl(r.URL.Path, "api/filterPosts")
	if errUrl == "not found" {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.NotFoundError)
		return
	}

	categories, err := models.ReadAllCategories(db)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	posts, err := models.FilterPosts(db, searchTerm)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		LoginUser  userManagementModels.User
		Posts      []models.Post
		Categories []models.Category
		SearchTerm string
	}{
		LoginUser:  userManagementModels.User{},
		Posts:      posts,
		Categories: categories,
		SearchTerm: searchTerm,
	}

	if loginStatus {
		data_obj_sender.LoginUser = loginUser
	}

	res := utils.Result{
		Success: true,
		Message: "Posts fetched successfully",
		Data:    data_obj_sender,
	}
	utils.ReturnJson(w, res)
}

func ReadMyCreatedPosts(w http.ResponseWriter, r *http.Request, db *dbTools.DBContainer) {
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

	categories, err := models.ReadAllCategories(db)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	posts, err := models.ReadPostsByUserId(db, loginUser.ID)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		LoginUser  userManagementModels.User
		Posts      []models.Post
		Categories []models.Category
	}{
		LoginUser:  loginUser,
		Posts:      posts,
		Categories: categories,
	}

	res := utils.Result{
		Success: true,
		Message: "Posts fetched successfully",
		Data:    data_obj_sender,
	}
	utils.ReturnJson(w, res)
}

func ReadMyLikedPosts(w http.ResponseWriter, r *http.Request, db *dbTools.DBContainer) {
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

	categories, err := models.ReadAllCategories(db)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	posts, err := models.ReadPostsLikedByUserId(db, loginUser.ID)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		LoginUser  userManagementModels.User
		Posts      []models.Post
		Categories []models.Category
	}{
		LoginUser:  loginUser,
		Posts:      posts,
		Categories: categories,
	}

	res := utils.Result{
		Success: true,
		Message: "Posts fetched successfully",
		Data:    data_obj_sender,
	}
	utils.ReturnJson(w, res)
}

func ReadPost(w http.ResponseWriter, r *http.Request, db *dbTools.DBContainer) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	loginStatus, loginUser, _, checkLoginError := userManagementControllers.CheckLogin(w, r, db)
	if checkLoginError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	uuid, errUrl := utils.ExtractUUIDFromUrl(r.URL.Path, "api/post")
	if errUrl == "not found" {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.NotFoundError)
		return
	}

	post, err := models.ReadPostByUUID(db, uuid, loginUser.ID)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		LoginUser userManagementModels.User
		Post      models.Post
		Comments  []models.Comment
	}{
		LoginUser: loginUser,
		Post:      post,
		Comments:  nil,
	}

	if loginStatus {
		comments, err := models.ReadAllCommentsForPostByUserID(db, post.ID, loginUser.ID)
		if err != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		data_obj_sender.Comments = comments
	} else {
		comments, err := models.ReadAllCommentsForPost(db, post.ID)
		if err != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		data_obj_sender.Comments = comments
	}

	// // Create a template with a function map
	// tmpl, err := template.New("post_details.html").Funcs(template.FuncMap{
	// 	"formatDate": utils.FormatDate, // Register function globally
	// }).ParseFiles(
	// 	publicUrl+"post_details.html",
	// 	publicUrl+"templates/header.html",
	// 	publicUrl+"templates/navbar.html",
	// 	publicUrl+"templates/footer.html",
	// )
	// if err != nil {
	// 	errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
	// 	return
	// }

	// // Execute template with data
	// err = tmpl.Execute(w, data_obj_sender)
	// if err != nil {
	// 	errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
	// }

	res := utils.Result{
		Success: true,
		Message: "Post submitted successfully",
		Data:    data_obj_sender,
	}
	utils.ReturnJson(w, res)
}

func SubmitPost(w http.ResponseWriter, r *http.Request, db *dbTools.DBContainer) {
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

	// Parse the multipart form with a max memory of 10MB
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
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
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}
		defer file.Close()

		// Call your file upload function
		uploadedFile, err := fileManagementControllers.FileUpload(file, handler)
		if err != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		uploadedFiles[handler.Filename] = uploadedFile
	}

	post := &models.Post{
		Title:   title,
		Content: content,
		UserId:  loginUser.ID,
	}

	// Convert the string slice to an int slice
	categoryIds := make([]int, 0, len(categories))
	for _, category := range categories {
		if id, err := strconv.Atoi(category); err == nil {
			categoryIds = append(categoryIds, id)
		} else {
			// Handle error if conversion fails (for example, invalid input)
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
			return
		}
	}

	// Insert a record while checking duplicates
	_, insertError := models.InsertPost(db, post, categoryIds, uploadedFiles)
	if insertError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	res := utils.Result{
		Success: true,
		Message: "Post submitted successfully",
		Data:    nil,
	}
	utils.ReturnJson(w, res)
}

func UpdatePost(w http.ResponseWriter, r *http.Request, db *dbTools.DBContainer) {
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

	// Parse the multipart form with a max memory of 10MB
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
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
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	// Retrieve all uploaded files
	files := r.MultipartForm.File["postFiles"]

	uploadedFiles := make(map[string]string)

	for _, handler := range files {
		file, err := handler.Open()
		if err != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}
		defer file.Close()

		// Call your file upload function
		uploadedFile, err := fileManagementControllers.FileUpload(file, handler)
		if err != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		uploadedFiles[handler.Filename] = uploadedFile
	}

	post := &models.Post{
		ID:      id,
		Title:   title,
		Content: content,
		UserId:  loginUser.ID,
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
	updateError := models.UpdatePost(db, post, categoryIds, uploadedFiles, loginUser.ID)
	if updateError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	res := utils.Result{
		Success: true,
		Message: "Post updated successfully",
		Data:    nil,
	}
	utils.ReturnJson(w, res)
}

func DeletePost(w http.ResponseWriter, r *http.Request, db *dbTools.DBContainer) {
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

	idStr := r.FormValue("id")

	if len(idStr) == 0 {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
		return
	}

	post_id, err := strconv.Atoi(idStr)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	// Update a record while checking duplicates
	updateError := models.UpdateStatusPost(db, post_id, "delete", loginUser.ID)
	if updateError != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	res := utils.Result{
		Success: true,
		Message: "Post removed successfully",
		Data:    nil,
	}
	utils.ReturnJson(w, res)
}

// func AdminDeletePost(w http.ResponseWriter, r *http.Request, db *dbTools.DBContainer) {
// 	if r.Method != http.MethodPost {
// 		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
// 		return
// 	}

// 	loginUser, ok := r.Context().Value(middlewares.AdminKey).(userManagementModels.User)
// 	if !ok {
// 		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
// 		return
// 	}

// 	err := r.ParseForm()
// 	if err != nil {
// 		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
// 		return
// 	}

// 	idStr := r.FormValue("id")

// 	if len(idStr) == 0 {
// 		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.BadRequestError)
// 		return
// 	}

// 	post_id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
// 		return
// 	}

// 	// Update a record while checking duplicates
// 	updateError := models.UpdateStatusPost(db, post_id, "delete", loginUser.ID)
// 	if updateError != nil {
// 		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
// 		return
// 	}

// 	//todo
// 	// userManagementControllers.RedirectToAdminIndex(w, r)
// }

func LikePost(w http.ResponseWriter, r *http.Request, db *dbTools.DBContainer) {
	if r.Method != http.MethodPost {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	// loginUser, ok := r.Context().Value(middlewares.UserContextKey).(userManagementModels.User)
	// if !ok {
	// 	errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.UnauthorizedError)
	// 	return
	// }
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

	existingLikeId, existingFeedbackRating := models.PostHasFeedbacked(db, loginUser.ID, postIDInt)

	var resMessage string
	if ratingStr == "1" {
		resMessage = "You liked successfully"
	} else {
		resMessage = "You disliked successfully"
	}

	rating, err := strconv.Atoi(ratingStr)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	if existingLikeId == -1 {
		post := &models.PostFeedback{
			Rating: rating,
			PostId: postIDInt,
			UserId: loginUser.ID,
		}
		_, insertError := models.InsertPostFeedback(db, post)
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
		updateError := models.UpdateStatusFeedback(db, existingLikeId, "delete", loginUser.ID)
		if updateError != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		if existingFeedbackRating != ratingStr { //this is duplicated like or duplicated dislike so we should update it to disable
			post := &models.PostFeedback{
				Rating: rating,
				PostId: postIDInt,
				UserId: loginUser.ID,
			}
			_, insertError := models.InsertPostFeedback(db, post)
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
