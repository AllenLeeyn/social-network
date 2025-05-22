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

	http.HandleFunc("/api/user/update",
		middleware.CheckHttpRequest("user", http.MethodPost,
			userContollers.UserUpdateHandler))

	// implement uuid as part of path
	http.HandleFunc("/api/users",
		middleware.CheckHttpRequest("user", http.MethodGet,
			userContollers.ViewUsersHandler))

	// implement uuid as part of path
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

	// implement uuid as part of path
	http.HandleFunc("/api/groups/",
		middleware.CheckHttpRequest("user", http.MethodGet,
			groupContollers.ViewGroupsHandler))

	// implement uuid as part of path
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

	// implement uuid as part of path
	http.HandleFunc("/api/followers/",
		middleware.CheckHttpRequest("user", http.MethodGet,
			followingContollers.ViewFollowingsHandler))

	// implement uuid as part of path
	http.HandleFunc("/api/follower/requests/",
		middleware.CheckHttpRequest("user", http.MethodGet,
			followingContollers.ViewFollowingRequestsHandler))
}
