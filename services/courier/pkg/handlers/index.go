package handlers

import (
	"diploma/services/courier/pkg/models"
	"diploma/services/courier/pkg/storage"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type OrderItemsData struct {
	OrderId int
	OrderItems []models.OrderItem
}

func CourierPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("front/pages/courier/courier.html")
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
	CheckCourierState(w, r, "http://localhost:8083/courier")

	tmpl, err := template.ParseFiles("front/pages/courier/view_orders.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Fatalf("StartPage: %s", err.Error())
	}

	orders, err := storage.ViewOrders()
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Fatalf("StartPage: %s", err.Error())
	}

	err = tmpl.Execute(w, orders)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalf("StartPage: %s", err.Error())
	}
}

func ViewOrderItemsPage(w http.ResponseWriter, r *http.Request) {
	CheckCourierState(w, r, "http://localhost:8083/courier")

	orderId, err := strconv.Atoi(r.FormValue("order_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("front/pages/courier/view_order_items.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Fatalf("StartPage: %s", err.Error())
	}

	orderItems, err := storage.ViewOrderItem(orderId)
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Fatalf("StartPage: %s", err.Error())
	}

	data := OrderItemsData{
		OrderId: orderId,
		OrderItems: orderItems,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalf("StartPage: %s", err.Error())
	}
}