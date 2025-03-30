package handlers

import (
	"diploma/services/customer/pkg/models"
	"diploma/services/customer/pkg/redis"
	"diploma/services/customer/pkg/storage"
	"encoding/json"

	"diploma/services/customer/pkg/mq"
	"html/template"

	"log"
	"strconv"

	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type ViewOrderItemsData struct {
	OrderId    int                 `json:"order_id"`
	OrderItems []models.BucketItem `json:"order_items"`
}

func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("customer")
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	tokenStr := c.Value
	redis.RedisClient.Del(tokenStr)

	http.SetCookie(w, &http.Cookie{
		Name:     "customer",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1,
	})

	http.Redirect(w, r, "http://localhost:8082/authorization/customer", http.StatusSeeOther)
}

func GetCustomerId(w http.ResponseWriter, r *http.Request) int {
	c, err := r.Cookie("customer")

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return -1
	}

	tokenStr := c.Value
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return -1
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return -1
	}

	if !token.Valid {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return -1
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return -1
	}

	customerId := claims["customer_id"].(float64)

	return int(customerId)
}

func MakeOrder(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	var bucketId int
	productsWithAmount := map[int]int{}
	for key, values := range r.Form {
		if len(values) == 0 {
			continue
		}

		value := values[0]

		if key == "bucket_id" {
			bucketId, err = strconv.Atoi(value)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else if len(key) > 7 && key[:7] == "amount_" {
			productId := key[7:]
			amount, err := strconv.Atoi(value)
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

	err = storage.UpdateBucketItems(bucketId, productsWithAmount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	customerId := GetCustomerId(w, r)
	err = storage.UpdateBucketStatus(bucketId, customerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	allProductCost, err := storage.GetAllProductCost(bucketId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = storage.MakeOrder(bucketId, customerId, allProductCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderId, created_at, err := storage.SelectOrderIdAndCreatedAt(bucketId, customerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderItems, err := storage.ViewOrderItems(bucketId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	madeOrderMessage := models.OrderMessage{
		OrderId:     orderId,
		CustomerId:  customerId,
		TotalCost:   float32(allProductCost),
		Status:      "created",
		CreatedAt:   created_at,
		OrderItems:  orderItems,
	}

	err = mq.ProduceMessage(madeOrderMessage, "Made order")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	MadeOrderPage(w, r)
}

func ViewOrderItems(w http.ResponseWriter, r *http.Request) {
	orderId, err := strconv.Atoi(r.FormValue("order_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bucketId, err := strconv.Atoi(r.FormValue("bucket_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderItems, err := storage.ViewOrderItems(bucketId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := ViewOrderItemsData{
		OrderId:    orderId,
		OrderItems: orderItems,
	}

	tmpl, err := template.ParseFiles("front/pages/customer/view_order_items.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Fatalf("StartPage: %s", err.Error())
	}

	if err = tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func GetOrderStatuses(w http.ResponseWriter, r *http.Request) {
	customerId := GetCustomerId(w, r)

	statuses, err := storage.GetOrderStatuses(customerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string][]models.OrderStatus{
		"statuses": statuses,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}