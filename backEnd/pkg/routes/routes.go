package routes

import (
	"net/http"

	chatContollers "social-network/pkg/chatManagement/controllers"
	groupContollers "social-network/pkg/groupManagement/controllers"
	middleware "social-network/pkg/middleware"
	userContollers "social-network/pkg/userManagement/controllers"
)

func SetupRoutes(cc *chatContollers.ChatController) {

	// ---------------------------- user management controller APIs ---------------------------- //
	http.HandleFunc("/api/register",
		middleware.CheckHttpRequest("guest", http.MethodPost, userContollers.RegisterHandler))

	http.HandleFunc("/api/login",
		middleware.CheckHttpRequest("guest", http.MethodPost, userContollers.LoginHandler))

	http.HandleFunc("/api/logout",
		middleware.CheckHttpRequest("user", http.MethodGet, userContollers.LogoutHandler))

	http.HandleFunc("/api/updateUser",
		middleware.CheckHttpRequest("user", http.MethodPost, userContollers.UpdateUserHandler))

	// to add
	http.HandleFunc("/api/users",
		middleware.CheckHttpRequest("user", http.MethodPost, userContollers.ViewUsersHandler))

	// to add
	http.HandleFunc("/api/user",
		middleware.CheckHttpRequest("user", http.MethodPost, userContollers.ViewUserHandler))

	http.HandleFunc("/api/ws",
		middleware.CheckHttpRequest("user", http.MethodGet, cc.WSHandler))

	http.HandleFunc("/api/createGroup",
		middleware.CheckHttpRequest("user", http.MethodPost, groupContollers.CreateGroupHandler))

	// to add
	http.HandleFunc("/api/groups",
		middleware.CheckHttpRequest("user", http.MethodPost, groupContollers.ViewGroupsHandler))

	// to add
	http.HandleFunc("/api/group",
		middleware.CheckHttpRequest("user", http.MethodPost, groupContollers.ViewGroupHandler))

	http.HandleFunc("/api/updateGroup",
		middleware.CheckHttpRequest("user", http.MethodPost, groupContollers.UpdateGroupHandler))

}
