package handlers

import (
	"diploma/services/customer/pkg/storage"
	"html/template"
	"log"
	"net/http"
)

type SuccessData struct {
	Title string
	Path  string
}

func CustomerPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("front/pages/customer/customer.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Fatalf("StartPage: %s", err.Error())
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalf("StartPage: %s", err.Error())
	}
}

func SelectCategoryPage(w http.ResponseWriter, r *http.Request) {
	categories, err := storage.ViewAllCategories()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("front/pages/customer/select_category.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = tmpl.Execute(w, categories); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
