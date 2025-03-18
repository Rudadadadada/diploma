package handlers

import (
	"log"
	"net/http"
	"html/template"
	"diploma/services/admin/pkg/storage"
)

type SuccessData struct {
    Title string
    Path  string
}

func AdminPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("front/pages/admin/admin.html")
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

func CategoriesPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("front/pages/admin/categories.html")
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

func CreateCategoryPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("front/pages/admin/create_category.html")
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

func RemoveCategoryPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("front/pages/admin/remove_category.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Fatalf("StartPage: %s", err.Error())
	}

	categories, err := storage.ViewAllCategories()

	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalf("StartPage: %s", err.Error())
	}

	err = tmpl.Execute(w, categories)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalf("StartPage: %s", err.Error())
	}
}

func ProductsPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("front/pages/admin/products.html")
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

func CreateProductPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("front/pages/admin/create_product.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Fatalf("StartPage: %s", err.Error())
	}

	categories, err := storage.ViewAllCategories()

	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalf("StartPage: %s", err.Error())
	}

	err = tmpl.Execute(w, categories)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalf("StartPage: %s", err.Error())
	}
}

func RemoveProductPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("front/pages/admin/remove_product.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Fatalf("StartPage: %s", err.Error())
	}

	products, _, err := storage.ViewAllProducts()

	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalf("StartPage: %s", err.Error())
	}

	err = tmpl.Execute(w, products)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalf("StartPage: %s", err.Error())
	}
}