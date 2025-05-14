package routes

import (
	"database/sql"
	"net/http"

	// socialMediaManagementControllers "social-network/pkg/socialMediaManagement/controllers"
	middleware "social-network/pkg/middleware"
	userContollers "social-network/pkg/userManagement/controllers"
)

func SetupRoutes(sqlDB *sql.DB) {

	// Create controller instance
	uc := userContollers.NewUserController(sqlDB)
	mw := middleware.SetUpMiddleware(sqlDB)

	// ---------------------------- user management controller APIs ---------------------------- //
	http.HandleFunc("/api/register",
		mw.CheckHttpRequest("guest", http.MethodPost, uc.RegisterHandler)) /*post method*/

	http.HandleFunc("/api/login",
		mw.CheckHttpRequest("guest", http.MethodPost, uc.LoginHandler)) /*post method*/

	http.HandleFunc("/api/logout/",
		mw.CheckHttpRequest("user", http.MethodGet, uc.LogoutHandler))

	http.HandleFunc("/api/updateUser",
		mw.CheckHttpRequest("user", http.MethodPost, uc.UpdateUserHandler)) /*post method*/
}
