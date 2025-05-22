package routes

import (
	"net/http"

	middleware "social-network/pkg/middleware"

	chatContollers "social-network/pkg/chatManagement/controllers"
	eventContollers "social-network/pkg/eventManagement/controllers"
	fileControllers "social-network/pkg/fileManagement/controllers"
	followingContollers "social-network/pkg/followingManagement/controllers"
	forumControllers "social-network/pkg/forumManagement/controllers"
	groupContollers "social-network/pkg/groupManagement/controllers"
	notificationControllers "social-network/pkg/notificationManagement/controllers"
	userContollers "social-network/pkg/userManagement/controllers"
	// socialMediaManagementControllers "social-network/pkg/socialMediaManagement/controllers"
)

func SetupRoutes(cc *chatContollers.ChatController) {

	// ---------------------------- user management controller APIs ---------------------------- //
	http.HandleFunc("/api/register",
		middleware.CheckHttpRequest("guest", http.MethodPost,
			userContollers.RegisterHandler))

	http.HandleFunc("/api/login",
		middleware.CheckHttpRequest("guest", http.MethodPost,
			userContollers.LoginHandler))

	http.HandleFunc("/api/logout",
		middleware.CheckHttpRequest("user", http.MethodGet,
			userContollers.LogoutHandler))

	http.HandleFunc("/api/user/update",
		middleware.CheckHttpRequest("user", http.MethodPost,
			userContollers.UserUpdateHandler))

	http.HandleFunc("/api/users",
		middleware.CheckHttpRequest("user", http.MethodGet,
			userContollers.ViewUsersHandler))

	http.HandleFunc("/api/user/",
		middleware.CheckHttpRequest("user", http.MethodGet,
			userContollers.ViewUserHandler))

	// ---------------------------- chat management controller APIs ---------------------------- //
	http.HandleFunc("/api/ws",
		middleware.CheckHttpRequest("user", http.MethodGet,
			cc.WSHandler))

	// ---------------------------- group management controller APIs ---------------------------- //
	http.HandleFunc("/api/group/create",
		middleware.CheckHttpRequest("user", http.MethodPost,
			groupContollers.GroupCreateHandler))

	http.HandleFunc("/api/group/update",
		middleware.CheckHttpRequest("user", http.MethodPost,
			groupContollers.GroupUpdateHandler))

	http.HandleFunc("/api/groups/",
		middleware.CheckHttpRequest("user", http.MethodGet,
			groupContollers.ViewGroupsHandler))

	http.HandleFunc("/api/group/",
		middleware.CheckHttpRequest("user", http.MethodGet,
			groupContollers.ViewGroupHandler))

	// ---------------------------- group member management controller APIs ---------------------------- //
	http.HandleFunc("/api/group/invite",
		middleware.CheckHttpRequest("user", http.MethodPost,
			groupContollers.GroupInviteRequestHandler))

	http.HandleFunc("/api/group/join",
		middleware.CheckHttpRequest("user", http.MethodPost,
			groupContollers.GroupJoinRequestHandler))

	http.HandleFunc("/api/group/quit",
		middleware.CheckHttpRequest("user", http.MethodPost,
			groupContollers.GroupQuitHandler))

	http.HandleFunc("/api/group/member/remove",
		middleware.CheckHttpRequest("user", http.MethodPost,
			groupContollers.GroupMemberRemoveHandler))

	http.HandleFunc("/api/group/member/response",
		middleware.CheckHttpRequest("user", http.MethodPost,
			groupContollers.GroupMemberResponseHandler))

	// implement uuid as part of path
	http.HandleFunc("/api/group/members/",
		middleware.CheckHttpRequest("user", http.MethodGet,
			groupContollers.ViewGroupMembersHandle))

	// implement uuid as part of path
	http.HandleFunc("/api/group/member/requests/",
		middleware.CheckHttpRequest("user", http.MethodGet,
			groupContollers.ViewGroupMemberRequestsHandle))

	// ---------------------------- group event management controller APIs ---------------------------- //
	http.HandleFunc("/api/group/event/create",
		middleware.CheckHttpRequest("user", http.MethodPost,
			eventContollers.EventCreateHandler))

	// ---------------------------- following management controller APIs ---------------------------- //
	http.HandleFunc("/api/follower/request",
		middleware.CheckHttpRequest("user", http.MethodPost,
			followingContollers.FollowingRequestHandler))

	http.HandleFunc("/api/follower/response",
		middleware.CheckHttpRequest("user", http.MethodPost,
			followingContollers.FollowingResponseHandler))

	http.HandleFunc("/api/follower/unfollow",
		middleware.CheckHttpRequest("user", http.MethodPost,
			followingContollers.UnfollowHandler))

	http.HandleFunc("/api/follower/remove",
		middleware.CheckHttpRequest("user", http.MethodPost,
			followingContollers.FollowingRemoveHandler))

	http.HandleFunc("/api/followers/",
		middleware.CheckHttpRequest("user", http.MethodGet,
			followingContollers.ViewFollowingsHandler))

	http.HandleFunc("/api/follower/requests/",
		middleware.CheckHttpRequest("user", http.MethodGet,
			followingContollers.ViewFollowingRequestsHandler))

	// ---------------------------- file management controller APIs ---------------------------- //
	http.HandleFunc("/api/uploadFile",
		middleware.CheckHttpRequest("user", http.MethodPost,
			fileControllers.FileUploadHandler)) /*post method*/

	// ---------------------------- forum management controller APIs ---------------------------- //
	http.HandleFunc("/api/categories/",
		middleware.CheckHttpRequest("user", http.MethodGet,
			forumControllers.ReadAllCategoriesHandler))

	http.HandleFunc("/api/allPosts/",
		middleware.CheckHttpRequest("user", http.MethodGet,
			forumControllers.ReadAllPostsHandler))

	http.HandleFunc("/api/myCreatedPosts/",
		middleware.CheckHttpRequest("user", http.MethodGet,
			forumControllers.ReadMyCreatedPostsHandler))

	http.HandleFunc("/api/myLikedPosts/",
		middleware.CheckHttpRequest("user", http.MethodGet,
			forumControllers.ReadMyLikedPostsHandler))

	// router.HandleFunc("/post/{id}", forumControllers.ReadPost).Methods("GET")
	http.HandleFunc("/api/post/",
		middleware.CheckHttpRequest("user", http.MethodGet,
			forumControllers.ReadPostHandler))

	// router.HandleFunc("/posts/{categoryName}", forumControllers.ReadPostsByCategory).Methods("GET")
	http.HandleFunc("/api/posts/",
		middleware.CheckHttpRequest("user", http.MethodGet,
			forumControllers.ReadPostsByCategoryHandler))

	http.HandleFunc("/api/filterPosts/",
		middleware.CheckHttpRequest("user", http.MethodGet,
			forumControllers.FilterPostsHandler))

	http.HandleFunc("/api/submitPost",
		middleware.CheckHttpRequest("user", http.MethodPost,
			forumControllers.SubmitPostHandler)) /*post method*/

	http.HandleFunc("/api/updatePost",
		middleware.CheckHttpRequest("user", http.MethodPut,
			forumControllers.UpdatePostHandler)) /*put method*/

	http.HandleFunc("/api/deletePost",
		middleware.CheckHttpRequest("user", http.MethodDelete,
			forumControllers.DeletePostHandler)) /*delete method*/

	http.HandleFunc("/api/postFeedback",
		middleware.CheckHttpRequest("user", http.MethodPost,
			forumControllers.PostFeedbackHandler)) /*post method*/

	http.HandleFunc("/api/submitComment",
		middleware.CheckHttpRequest("user", http.MethodPost,
			forumControllers.SubmitCommentHandler)) /*post method*/

	http.HandleFunc("/api/updateComment",
		middleware.CheckHttpRequest("user", http.MethodPut,
			forumControllers.UpdateCommentHandler)) /*put method*/

	http.HandleFunc("/api/deleteComment",
		middleware.CheckHttpRequest("user", http.MethodDelete,
			forumControllers.DeleteCommentHandler)) /*delete method*/

	http.HandleFunc("/api/commentFeedback",
		middleware.CheckHttpRequest("user", http.MethodPost,
			forumControllers.CommentFeedbackHandler)) /*post method*/

	// ---------------------------- notification management controller APIs ---------------------------- //
	http.HandleFunc("/api/submitNotification",
		middleware.CheckHttpRequest("user", http.MethodPost,
			notificationControllers.SubmitNotificationHandler)) /*post method*/

	http.HandleFunc("/api/updateNotificationReadStatus",
		middleware.CheckHttpRequest("user", http.MethodPost,
			notificationControllers.UpdateNotificationReadStatusHandler)) /*post method*/

	http.HandleFunc("/api/deleteNotification",
		middleware.CheckHttpRequest("user", http.MethodPost,
			notificationControllers.DeleteNotificationHandler)) /*post method*/

	http.HandleFunc("/api/notifications/",
		middleware.CheckHttpRequest("user", http.MethodGet,
			notificationControllers.ReadAllNotificationsHandler))

	// router.HandleFunc("/notification/{id}", notificationControllers.ReadNotificationByIdHandler).Methods("GET")
	http.HandleFunc("/api/notification/",
		middleware.CheckHttpRequest("user", http.MethodGet,
			notificationControllers.ReadNotificationByIdHandler))
}
