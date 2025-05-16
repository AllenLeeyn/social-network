package routes

import (
	"database/sql"
	"net/http"

	// socialMediaManagementControllers "social-network/pkg/socialMediaManagement/controllers"
	fileControllers "social-network/pkg/fileManagement/controllers"
	forumControllers "social-network/pkg/forumManagement/controllers"
	middleware "social-network/pkg/middleware"
	userContollers "social-network/pkg/userManagement/controllers"
)

func SetupRoutes(sqlDB *sql.DB) {
	// ---------------------------- user management controller APIs ---------------------------- //
	http.HandleFunc("/api/register",
		middleware.CheckHttpRequest("guest", http.MethodPost, userContollers.RegisterHandler)) /*post method*/

	http.HandleFunc("/api/login",
		middleware.CheckHttpRequest("guest", http.MethodPost, userContollers.LoginHandler)) /*post method*/

	http.HandleFunc("/api/logout/",
		middleware.CheckHttpRequest("user", http.MethodGet, userContollers.LogoutHandler))

	http.HandleFunc("/api/updateUser",
		middleware.CheckHttpRequest("user", http.MethodPost, userContollers.UpdateUserHandler)) /*post method*/

	http.HandleFunc("/api/uploadFile",
		middleware.CheckHttpRequest("user", http.MethodPost, fileControllers.FileUploadHandler)) /*post method*/

	http.HandleFunc("/api/categories/",
		middleware.CheckHttpRequest("user", http.MethodGet, forumControllers.ReadAllCategoriesHandler))

	http.HandleFunc("/api/allPosts/",
		middleware.CheckHttpRequest("user", http.MethodGet, forumControllers.ReadAllPostsHandler)) //todo read only my visible posts

	http.HandleFunc("/api/myCreatedPosts/",
		middleware.CheckHttpRequest("user", http.MethodGet, forumControllers.ReadMyCreatedPostsHandler))

	http.HandleFunc("/api/myLikedPosts/",
		middleware.CheckHttpRequest("user", http.MethodGet, forumControllers.ReadMyLikedPostsHandler))

	// router.HandleFunc("/post/{id}", forumControllers.ReadPost).Methods("GET")
	http.HandleFunc("/api/post/",
		middleware.CheckHttpRequest("user", http.MethodGet, forumControllers.ReadPostHandler))

	// router.HandleFunc("/posts/{categoryName}", forumControllers.ReadPostsByCategory).Methods("GET")
	http.HandleFunc("/api/posts/",
		middleware.CheckHttpRequest("user", http.MethodGet, forumControllers.ReadPostsByCategoryHandler))

	http.HandleFunc("/api/filterPosts/",
		middleware.CheckHttpRequest("user", http.MethodGet, forumControllers.FilterPostsHandler))

	http.HandleFunc("/api/submitPost",
		middleware.CheckHttpRequest("user", http.MethodPost, forumControllers.SubmitPostHandler)) /*post method*/ //todo fill post audiences

	http.HandleFunc("/api/updatePost",
		middleware.CheckHttpRequest("user", http.MethodPut, forumControllers.UpdatePostHandler)) /*put method*/ //todo search for duplicate handling while updating categories

	http.HandleFunc("/api/deletePost",
		middleware.CheckHttpRequest("user", http.MethodDelete, forumControllers.DeletePostHandler)) /*delete method*/

	http.HandleFunc("/api/postFeedback",
		middleware.CheckHttpRequest("user", http.MethodPost, forumControllers.PostFeedbackHandler)) /*post method*/

	http.HandleFunc("/api/submitComment",
		middleware.CheckHttpRequest("user", http.MethodPost, forumControllers.SubmitCommentHandler)) /*post method*/

	http.HandleFunc("/api/updateComment",
		middleware.CheckHttpRequest("user", http.MethodPut, forumControllers.UpdateCommentHandler)) /*post method*/

	http.HandleFunc("/api/deleteComment",
		middleware.CheckHttpRequest("user", http.MethodDelete, forumControllers.DeleteCommentHandler)) /*post method*/

	http.HandleFunc("/api/commentFeedback",
		middleware.CheckHttpRequest("user", http.MethodPost, forumControllers.CommentFeedbackHandler)) /*post method*/
}
