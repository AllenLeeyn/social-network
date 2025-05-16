package routes

import (
	"database/sql"
	"net/http"

	chatContollers "social-network/pkg/chatManagement/controllers"
	groupContollers "social-network/pkg/groupManagement/controllers"
	middleware "social-network/pkg/middleware"
	userContollers "social-network/pkg/userManagement/controllers"
)

func SetupRoutes(sqlDB *sql.DB, cc *chatContollers.ChatController) {

	// ---------------------------- user management controller APIs ---------------------------- //
	http.HandleFunc("/api/register",
		middleware.CheckHttpRequest("guest", http.MethodPost, userContollers.RegisterHandler))

	http.HandleFunc("/api/login",
		middleware.CheckHttpRequest("guest", http.MethodPost, userContollers.LoginHandler))

	http.HandleFunc("/api/logout",
		middleware.CheckHttpRequest("user", http.MethodGet, userContollers.LogoutHandler))

	http.HandleFunc("/api/updateUser",
		middleware.CheckHttpRequest("user", http.MethodPost, userContollers.UpdateUserHandler))

	http.HandleFunc("/api/ws",
		middleware.CheckHttpRequest("user", http.MethodGet, cc.WSHandler))

	http.HandleFunc("/api/createGroup",
		middleware.CheckHttpRequest("user", http.MethodPost, groupContollers.CreateGroupHandler))

	http.HandleFunc("/api/updateGroup",
		middleware.CheckHttpRequest("user", http.MethodPost, groupContollers.UpdateGroupHandler))

}
