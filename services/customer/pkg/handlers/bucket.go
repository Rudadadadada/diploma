package handlers

import (
	"diploma/services/customer/pkg/storage"
	"encoding/json"
	"strconv"

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

	productsWithAmount := map[int]int{}
	for key, values := range r.Form {
		if keyPrefix := "product_"; len(key) > len(keyPrefix) && key[:len(keyPrefix)] == keyPrefix {
			productId := key[len(keyPrefix):]

			if values[0] == "true" {
				amountKey := "amount_" + productId
				amountValues, amountExists := r.Form[amountKey]

				if amountExists && len(amountValues) > 0 {
					amount, err := strconv.Atoi(amountValues[0])
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					intProductId, err := strconv.Atoi(productId)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					productsWithAmount[intProductId] = amount
				}
			}
		}
	}

	customerId := GetCustomerId(w, r)

	err = storage.InertIntoBucket(customerId, productsWithAmount)

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
