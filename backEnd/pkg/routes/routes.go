package routes

import (
	"net/http"

	chatContollers "social-network/pkg/chatManagement/controllers"
	followingContollers "social-network/pkg/followingManagement/controllers"
	groupContollers "social-network/pkg/groupManagement/controllers"
	middleware "social-network/pkg/middleware"
	userContollers "social-network/pkg/userManagement/controllers"
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

	http.HandleFunc("/api/updateUser",
		middleware.CheckHttpRequest("user", http.MethodPost,
			userContollers.UpdateUserHandler))

	http.HandleFunc("/api/users",
		middleware.CheckHttpRequest("user", http.MethodGet,
			userContollers.ViewUsersHandler))

	http.HandleFunc("/api/user",
		middleware.CheckHttpRequest("user", http.MethodGet,
			userContollers.ViewUserHandler))

	http.HandleFunc("/api/ws",
		middleware.CheckHttpRequest("user", http.MethodGet,
			cc.WSHandler))

	http.HandleFunc("/api/createGroup",
		middleware.CheckHttpRequest("user", http.MethodPost,
			groupContollers.CreateGroupHandler))

	http.HandleFunc("/api/groups",
		middleware.CheckHttpRequest("user", http.MethodGet,
			groupContollers.ViewGroupsHandler))

	http.HandleFunc("/api/group",
		middleware.CheckHttpRequest("user", http.MethodGet,
			groupContollers.ViewGroupHandler))

	http.HandleFunc("/api/updateGroup",
		middleware.CheckHttpRequest("user", http.MethodPost,
			groupContollers.UpdateGroupHandler))

	http.HandleFunc("/api/followRequest",
		middleware.CheckHttpRequest("user", http.MethodPost,
			followingContollers.FollowingRequestHandler))

	http.HandleFunc("/api/followResponse",
		middleware.CheckHttpRequest("user", http.MethodPost,
			followingContollers.FollowingResponseHandler))

	http.HandleFunc("/api/unfollowRequest",
		middleware.CheckHttpRequest("user", http.MethodPost,
			followingContollers.UnfollowHandler))

	http.HandleFunc("/api/removeFollower",
		middleware.CheckHttpRequest("user", http.MethodPost,
			followingContollers.RemoveFollowerHandler))

	http.HandleFunc("/api/followingRequests",
		middleware.CheckHttpRequest("user", http.MethodGet,
			followingContollers.ViewFollowingRequestsHandler))

	// to return data if user is follower or profile is public
	http.HandleFunc("/api/followings",
		middleware.CheckHttpRequest("user", http.MethodGet,
			followingContollers.ViewFollowingsHandler))
}
