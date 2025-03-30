package handlers

import (
	"diploma/services/customer/pkg/models"
	"diploma/services/customer/pkg/storage"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type SuccessData struct {
	Title string
	Path  string
}

type BucketData struct {
	BucketItems     []models.BucketItem
	BucketId        int
	AllProductsCost float64
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

func SelectProductsByCategoryIdPage(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.Atoi(r.URL.Query().Get("category_id"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	products, err := storage.ViewProductsByCategoryId(categoryID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("front/pages/customer/select_products_by_category_id.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = tmpl.Execute(w, products); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func InsertedIntoBucketPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("front/pages/customer/inserted_into_bucket.html")
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

func BucketPage(w http.ResponseWriter, r *http.Request) {
	customerId := GetCustomerId(w, r)
	bucketItems, bucketId, err := storage.ViewBucket(customerId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("front/pages/customer/bucket.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var totalCost float32 = 0
	for _, tmp := range bucketItems {
		totalCost += tmp.TotalCost
	}

	data := BucketData{
		BucketItems:     bucketItems,
		BucketId:        bucketId,
		AllProductsCost: float64(totalCost),
	}

	if err = tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func MadeOrderPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("front/pages/customer/made_order.html")
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

func ViewOrdersPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("front/pages/customer/view_orders.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Fatalf("StartPage: %s", err.Error())
	}

	customerId := GetCustomerId(w, r)
	orders, err := storage.ViewOrders(customerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = tmpl.Execute(w, orders); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}