package handlers

import (
	"diploma/services/customer/pkg/storage"
	"encoding/json"
	"strconv"
	"strings"

	// "diploma/services/customer/pkg/storage"
	// "html/template"
	"net/http"
	// "strconv"
)

type RemoveItemRequest struct {
	BucketId  int `json:"bucket_id"`
	ProductId int `json:"product_id"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func InsertIntoBucket(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	productsStr := r.FormValue("selectedProducts")
	if productsStr == "" {
		http.Error(w, "No products selected", http.StatusBadRequest)
		return
	}

	productIds := strings.Split(productsStr, ",")
	productsWithAmount := map[int]int{}
	for _, productIdStr := range productIds {
		productId, err := strconv.Atoi(productIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		amountKey := "amount_" + productIdStr
		amountValues, amountExists := r.Form[amountKey]
		if amountExists && len(amountValues) > 0 {
			amount, err := strconv.Atoi(amountValues[0])
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			productsWithAmount[productId] = amount
		}
	}

	customerId := GetCustomerId(w, r)

	err = storage.InsertIntoBucket(customerId, productsWithAmount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	InsertedIntoBucketPage(w, r)
}


func RemoveItemFromBucket(w http.ResponseWriter, r *http.Request) {
	var req RemoveItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	err := storage.RemoveItemFromBucket(req.BucketId, req.ProductId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Success: true})
}
