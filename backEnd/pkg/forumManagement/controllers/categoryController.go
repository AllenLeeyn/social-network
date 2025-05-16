package controller

import (
	"net/http"
	errorControllers "social-network/pkg/errorManagement/controllers"
	"social-network/pkg/forumManagement/models"
	"social-network/pkg/utils"

	_ "github.com/mattn/go-sqlite3"
)

func ReadAllCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := models.ReadAllCategories()
	if err != nil {
		errorControllers.ErrorHandler(w, r, errorControllers.InternalServerError)
		return
	}

	utils.ReturnJsonSuccess(w, "Categories fetched successfully", categories)
}
