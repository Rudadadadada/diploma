package handlers

import (
	"diploma/services/admin/pkg/models"
	"diploma/services/admin/pkg/storage"
	"encoding/json"
	"log"
	"net/http"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Print("create product")
	json.NewEncoder(w).Encode(product)
}

func EditProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Print("edit product")
	json.NewEncoder(w).Encode(product)
}

func ViewAllProducts(w http.ResponseWriter, r *http.Request) {
	products, _ := storage.ViewAllProducts()

	log.Print(products)
	json.NewEncoder(w).Encode(products)
}