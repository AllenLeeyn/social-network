package routes

import (
	"database/sql"
	"net/http"

	// socialMediaManagementControllers "social-network/pkg/socialMediaManagement/controllers"
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
		middleware.CheckHttpRequest("user", http.MethodGet, userContollers.Logout))

	http.HandleFunc("/api/updateUser",
		middleware.CheckHttpRequest("user", http.MethodPost, userContollers.UpdateUser)) /*post method*/

	http.HandleFunc("/api/categories/",
		middleware.CheckHttpRequest("user", http.MethodGet, forumControllers.ReadAllCategories))

	http.HandleFunc("/api/allPosts/",
		middleware.CheckHttpRequest("user", http.MethodGet, forumControllers.ReadAllPosts))

	http.HandleFunc("/api/submitPost",
		middleware.CheckHttpRequest("user", http.MethodPost, forumControllers.SubmitPost)) /*post method*/

	http.HandleFunc("/api/myCreatedPosts/",
		middleware.CheckHttpRequest("user", http.MethodGet, forumControllers.ReadMyCreatedPosts))

	http.HandleFunc("/api/myLikedPosts/",
		middleware.CheckHttpRequest("user", http.MethodGet, forumControllers.ReadMyLikedPosts))

	// router.HandleFunc("/post/{id}", forumControllers.ReadPost).Methods("GET")
	http.HandleFunc("/api/post/",
		middleware.CheckHttpRequest("user", http.MethodGet, forumControllers.ReadPost))

	// router.HandleFunc("/posts/{categoryName}", forumControllers.ReadPostsByCategory).Methods("GET")
	http.HandleFunc("/api/posts/",
		middleware.CheckHttpRequest("user", http.MethodGet, forumControllers.ReadPostsByCategory))

	http.HandleFunc("/api/filterPosts/",
		middleware.CheckHttpRequest("user", http.MethodGet, forumControllers.FilterPosts))

	http.HandleFunc("/api/postFeedback",
		middleware.CheckHttpRequest("user", http.MethodPost, forumControllers.PostFeedback)) /*post method*/

	http.HandleFunc("/api/updatePost",
		middleware.CheckHttpRequest("user", http.MethodPost, forumControllers.UpdatePost)) /*post method*/

	http.HandleFunc("/api/deletePost",
		middleware.CheckHttpRequest("user", http.MethodPost, forumControllers.DeletePost)) /*post method*/

	http.HandleFunc("/api/likeComment",
		middleware.CheckHttpRequest("user", http.MethodPost, forumControllers.FeedbackComment)) /*post method*/

	http.HandleFunc("/api/submitComment",
		middleware.CheckHttpRequest("user", http.MethodPost, forumControllers.SubmitComment)) /*post method*/

	http.HandleFunc("/api/updateComment",
		middleware.CheckHttpRequest("user", http.MethodPost, forumControllers.UpdateComment)) /*post method*/

	http.HandleFunc("/api/deleteComment",
		middleware.CheckHttpRequest("user", http.MethodPost, forumControllers.DeleteComment)) /*post method*/
}
