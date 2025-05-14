package controller

import (
	"net/http"
	errorControllers "social-network/pkg/errorManagement/controllers"
	"social-network/pkg/forumManagement/models"
	"social-network/pkg/utils"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

const publicUrl = "modules/forumManagement/views/"

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Parse template
	tmpl, err := template.ParseFiles(publicUrl + "index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// If the URL is not exactly "/", respond with 404
		errorControllers.HandleErrorPage(w, r, errorControllers.NotFoundError)
		return
	}

	userIDRaw := r.Context().Value("userID")
	userID, isOk := userIDRaw.(int)
	if !isOk {
		errorControllers.HandleErrorPage(w, r, errorControllers.InternalServerError)
		return
	}

	categories, err := models.ReadAllCategories()
	if err != nil {
		errorControllers.HandleErrorPage(w, r, errorControllers.InternalServerError)
		return
	}

	posts, err := models.ReadAllPosts(userID)
	if err != nil {
		errorControllers.HandleErrorPage(w, r, errorControllers.InternalServerError)
		return
	}

	data_obj_sender := struct {
		Posts      []models.Post
		Categories []models.Category
	}{
		Posts:      posts,
		Categories: categories,
	}

	// Create a template with a function map
	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"formatDate": utils.FormatDate, // Register function globally
	}).ParseFiles(
		publicUrl + "index.html",
	)
	if err != nil {
		errorControllers.HandleErrorPage(w, r, errorControllers.InternalServerError)
		return
	}

	err = tmpl.Execute(w, data_obj_sender)
	if err != nil {
		errorControllers.HandleErrorPage(w, r, errorControllers.InternalServerError)
		return
	}

}
