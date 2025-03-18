package handlers

import (
	"diploma/services/admin/pkg/models"
	"diploma/services/admin/pkg/storage"
	"html/template"
	"net/http"
	"strconv"
)


func ViewAllProducts(w http.ResponseWriter, r *http.Request) {
	products, categoriesNames, err := storage.ViewAllProducts()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("front/pages/admin/view_products.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	data := []struct{
		Id          uint
		Name        string
		Amount      uint
		Cost        float32
		CategoryName string
	}{}
	
	for i := 0; i < len(products); i++ {
		tmp := struct {
			Id          uint
			Name        string
			Amount      uint
			Cost        float32
			CategoryName string
		}{
			Id:           products[i].Id,
			Name:         products[i].Name,
			Amount:       products[i].Amount,
			Cost:         products[i].Cost,
			CategoryName: categoriesNames[i],
		}

		data = append(data, tmp)
	}

	if err = tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("Product name")

	cost, err := strconv.Atoi(r.FormValue("Amount"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("Cost"), 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	categoryId, err := strconv.Atoi(r.FormValue("category"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newProduct := models.Product{
		Name: name,
		Amount: uint(cost),
		Cost: float32(amount),
		CategoryID: uint(categoryId),
	}	
	
	err = storage.CreateProduct(newProduct)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("front/pages/admin/success.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := SuccessData{
		Title: "Продукт успешно создан",
		Path: "/admin/products",
	}

	if err = tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func RemoveProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.FormValue("product"), 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	RemoveProduct := models.Product{Id: uint(id)}

	err = storage.RemoveProduct(RemoveProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("front/pages/admin/success.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := SuccessData{
		Title: "Продукт успешно удален",
		Path: "/admin/products",
	}

	if err = tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}