package handlers

import (
	"diploma/services/admin/pkg/models"
	"encoding/json"
	"log"
	"net/http"
)

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Print("create category")
	json.NewEncoder(w).Encode(category)
}

func EditCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Print("edit category")
	json.NewEncoder(w).Encode(category)
}
