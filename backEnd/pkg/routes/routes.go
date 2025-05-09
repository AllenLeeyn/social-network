package routes

import (
	"database/sql"
	"net/http"

	// socialMediaManagementControllers "social-network/pkg/socialMediaManagement/controllers"
	userManagementControllers "social-network/pkg/userManagement/controllers"
)

func SetupRoutes(sqlDB *sql.DB) {
	// ---------------------------- user management controller APIs ---------------------------- //
	http.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		userManagementControllers.RegisterHandler(w, r, sqlDB)
	}) /*post method*/
	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		userManagementControllers.LoginHandler(w, r, sqlDB)
	}) /*post method*/
	http.HandleFunc("/api/logout/", func(w http.ResponseWriter, r *http.Request) {
		userManagementControllers.Logout(w, r, sqlDB)
	})
	http.HandleFunc("/api/updateUser", func(w http.ResponseWriter, r *http.Request) {
		userManagementControllers.UpdateUser(w, r, sqlDB)
	}) /*post method*/
}
