package routes

import (
	"database/sql"
	"net/http"

	chatContollers "social-network/pkg/chatManagement/controllers"
	middleware "social-network/pkg/middleware"
	userContollers "social-network/pkg/userManagement/controllers"
)

func SetupRoutes(sqlDB *sql.DB, cc *chatContollers.ChatController) {

	// ---------------------------- user management controller APIs ---------------------------- //
	http.HandleFunc("/api/register",
		middleware.CheckHttpRequest("guest", http.MethodPost, userContollers.RegisterHandler)) /*post method*/

	http.HandleFunc("/api/login",
		middleware.CheckHttpRequest("guest", http.MethodPost, userContollers.LoginHandler)) /*post method*/

	http.HandleFunc("/api/logout",
		middleware.CheckHttpRequest("user", http.MethodGet, userContollers.LogoutHandler))

	http.HandleFunc("/api/updateUser",
		middleware.CheckHttpRequest("user", http.MethodPost, userContollers.UpdateUserHandler)) /*post method*/

	http.HandleFunc("/api/ws",
		middleware.CheckHttpRequest("user", http.MethodGet, cc.WSHandler)) /*post method*/

}
