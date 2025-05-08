package routes

import (
	"net/http"
	"social-network/pkg/dbTools"

	// socialMediaManagementControllers "social-network/pkg/socialMediaManagement/controllers"
	userManagementControllers "social-network/pkg/userManagement/controllers"
)

func SetupRoutes(db *dbTools.DBContainer) {
	// ---------------------------- user management controller APIs ---------------------------- //
	http.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		userManagementControllers.RegisterHandler(w, r, db)
	}) /*post method*/
	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		userManagementControllers.LoginHandler(w, r, db)
	}) /*post method*/
	http.HandleFunc("/api/logout/", func(w http.ResponseWriter, r *http.Request) {
		userManagementControllers.Logout(w, r, db)
	})
	http.HandleFunc("/api/updateUser", func(w http.ResponseWriter, r *http.Request) {
		userManagementControllers.UpdateUser(w, r, db)
	}) /*post method*/
	http.HandleFunc("/api/check-session", func(w http.ResponseWriter, r *http.Request) {
		userManagementControllers.CheckSessionHandler(w, r, db)
	})
}
