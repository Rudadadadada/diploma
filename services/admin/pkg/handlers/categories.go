package handlers

import (
	"diploma/services/admin/pkg/models"
	"diploma/services/admin/pkg/mq"
	"diploma/services/admin/pkg/storage"
	"html/template"
	"net/http"
	"strconv"
)

func ViewAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := storage.ViewAllCategories()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("front/pages/admin/view_categories.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err = tmpl.Execute(w, categories); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("Category name")
	newCategory := models.Category{Name: name}

	err := storage.CreateCategory(newCategory)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("front/pages/admin/success.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	syncDatabasesMessage, err := storage.SyncDatabases()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = mq.ProduceSyncMessage(*syncDatabasesMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := SuccessData{
		Title: "Категория успешно создана",
		Path: "/admin/categories",
	}

	if err = tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func RemoveCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.FormValue("category"), 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	removeCategory := models.Category{Id: uint(id)}

	err = storage.RemoveCategory(removeCategory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("front/pages/admin/success.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	syncDatabasesMessage, err := storage.SyncDatabases()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = mq.ProduceSyncMessage(*syncDatabasesMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := SuccessData{
		Title: "Категория успешно удалена",
		Path: "/admin/categories",
	}
	
	if err = tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}