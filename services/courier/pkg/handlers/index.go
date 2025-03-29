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
	OrderId    int
	OrderItems []models.OrderItem
	AllProductsCost float32
}

func CourierPage(w http.ResponseWriter, r *http.Request) {
	path, err := CheckCourierInProgress(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if path != nil {
		http.Redirect(w, r, *path, http.StatusSeeOther)
		return
	}

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
	path, err := CheckCourierInProgress(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if path != nil {
		http.Redirect(w, r, *path, http.StatusSeeOther)
		return
	}
	CheckCourierActive(w, r)

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
	path, err := CheckCourierInProgress(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if path != nil {
		http.Redirect(w, r, *path, http.StatusSeeOther)
		return
	}
	CheckCourierActive(w, r)

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
		OrderId:    orderId,
		OrderItems: orderItems,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalf("StartPage: %s", err.Error())
	}
}

func InProgressPage(w http.ResponseWriter, r *http.Request) {
	cameFrom := "in progress"
	path, err := CheckCourierInProgress(w, r, &cameFrom)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if path != nil {
		http.Redirect(w, r, *path, http.StatusSeeOther)
		return
	}

	CheckCourierActive(w, r)

	courierId := GetCourierId(w, r)

	orderId, err := storage.GetOrderId(courierId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("front/pages/courier/in_progress.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Fatalf("StartPage: %s", err.Error())
	}

	orderItems, err := storage.ViewOrderItem(orderId)
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Fatalf("StartPage: %s", err.Error())
	}

	allProductsCost, err := storage.GetOrderCost(orderId)
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Fatalf("StartPage: %s", err.Error())
	}

	data := OrderItemsData{
		OrderId:    orderId,
		OrderItems: orderItems,
		AllProductsCost: allProductsCost,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalf("StartPage: %s", err.Error())
	}
}
