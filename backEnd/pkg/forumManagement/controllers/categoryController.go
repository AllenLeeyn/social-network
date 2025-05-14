package controller

import (
	"database/sql"
	"net/http"
	errorManagementControllers "social-network/pkg/errorManagement/controllers"
	"social-network/pkg/forumManagement/models"
	"social-network/pkg/utils"

	_ "github.com/mattn/go-sqlite3"
)

func ReadAllCategories(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	categories, err := models.ReadAllCategories(db)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	res := utils.Result{
		Success: true,
		Message: "Post submitted successfully",
		Data:    categories,
	}
	utils.ReturnJson(w, res)
}
