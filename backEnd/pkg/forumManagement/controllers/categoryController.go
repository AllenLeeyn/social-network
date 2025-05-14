package controller

import (
	"net/http"
	errorControllers "social-network/pkg/errorManagement/controllers"
	"social-network/pkg/forumManagement/models"
	"social-network/pkg/utils"

	_ "github.com/mattn/go-sqlite3"
)

func ReadAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := models.ReadAllCategories()
	if err != nil {
		errorControllers.HandleErrorPage(w, r, errorControllers.InternalServerError)
		return
	}

	res := utils.Result{
		Success: true,
		Message: "Post submitted successfully",
		Data:    categories,
	}
	utils.ReturnJson(w, res)
}
