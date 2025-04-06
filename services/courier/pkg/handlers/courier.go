package handlers

import (
	"diploma/services/courier/pkg/models"
	"diploma/services/courier/pkg/mq"
	"diploma/services/courier/pkg/redis"
	"diploma/services/courier/pkg/storage"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("courier")
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	tokenStr := c.Value
	redis.RedisClient.Del(tokenStr)

	http.SetCookie(w, &http.Cookie{
		Name:     "courier",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1,
	})

	http.Redirect(w, r, "http://localhost:8082/authorization/courier", http.StatusSeeOther)
}

func GetCourierId(w http.ResponseWriter, r *http.Request) int {
	c, err := r.Cookie("courier")

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

	courierId := claims["courier_id"].(float64)

	return int(courierId)
}

func SetState(w http.ResponseWriter, r *http.Request) {
	path, err := CheckCourierInProgress(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if path != nil {
		http.Redirect(w, r, *path, http.StatusSeeOther)
		return
	}

	courierId := GetCourierId(w, r)

	err = storage.SetActive(courierId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "http://localhost:8083/courier", http.StatusSeeOther)
}

func GetState(w http.ResponseWriter, r *http.Request) {
	courierId := GetCourierId(w, r)

	err := storage.AddCourier(courierId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ptr, err := storage.GetState(courierId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var courier models.Courier
	if ptr != nil {
		courier = *ptr
	}

	response := map[string]bool{"active": courier.Active}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Ошибка при форматировании ответа", http.StatusInternalServerError)
		return
	}

	err = mq.ProduceState(courier, "Courier state")
	if err != nil {
		http.Error(w, "Ошибка при форматировании ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func CheckCourierActive(w http.ResponseWriter, r *http.Request) {
	courierId := GetCourierId(w, r)

	ptr, err := storage.GetState(courierId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var courier models.Courier
	if ptr != nil {
		courier = *ptr
	}

	if !courier.Active {
		http.Redirect(w, r, "http://localhost:8083/courier", http.StatusSeeOther)
		return
	}
}

func CheckCourierInProgress(w http.ResponseWriter, r *http.Request, cameFrom *string) (*string, error) {
	courierId := GetCourierId(w, r)

	ptr, err := storage.GetState(courierId)
	if err != nil {
		return nil, err
	}

	var courier models.Courier
	if ptr != nil {
		courier = *ptr
	}

	var path string
	if courier.InProgress {
		if cameFrom != nil && (*cameFrom == "in progress" ||
			*cameFrom == "take order from shop" ||
			*cameFrom == "get order status" ||
			*cameFrom == "not yet" ||
			*cameFrom == "declined" ||
			*cameFrom == "finish delivery") {
			return nil, err
		}

		path = "http://localhost:8083/courier/in_progress"
		return &path, nil

	} else {
		if cameFrom != nil && (*cameFrom == "in progress" ||
			*cameFrom == "take order from shop" ||
			*cameFrom == "get order status" ||
			*cameFrom == "not yet" ||
			*cameFrom == "declined" ||
			*cameFrom == "finish delivery") {
			path = "http://localhost:8083/courier"
			return &path, nil
		}
	}

	return nil, nil
}

func TakeOrder(w http.ResponseWriter, r *http.Request) {
	CheckCourierActive(w, r)

	orderId, err := strconv.Atoi(r.FormValue("order_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, ptr, err := storage.CheckOrderTaken(orderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if id == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status, err := storage.GetOrderStatus(orderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if status == "order declined" {
		http.Redirect(w, r, "http://localhost:8083/order/declined", http.StatusSeeOther)
		return
	}

	var took bool
	if ptr != nil {
		took = *ptr
	}

	if took {
		tmpl, err := template.ParseFiles("front/pages/courier/taken_by_other_courier.html")
		if err != nil {
			http.Error(w, err.Error(), 400)
			log.Fatalf("StartPage: %s", err.Error())
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), 500)
			log.Fatalf("StartPage: %s", err.Error())
		}
		return
	}

	courierId := GetCourierId(w, r)
	err = storage.TakeOrder(orderId, courierId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = storage.SetInProgress(courierId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ptrOrderMessage, err := storage.GetFullOrderInfo(orderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var tookOrderMessage models.OrderMessage
	if ptrOrderMessage != nil {
		tookOrderMessage = *ptrOrderMessage
	}

	ptrCourier, err := storage.GetState(courierId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var courier models.Courier
	if ptrCourier != nil {
		courier = *ptrCourier
	}

	tookOrderMessage.DeliveryStartedAt = time.Now()
	tookOrderMessage.Courier = courier

	err = mq.ProduceMessage(tookOrderMessage, "Order taken")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "http://localhost:8083/courier/in_progress", http.StatusSeeOther)
}

func GetOrderStatus(w http.ResponseWriter, r *http.Request) {
	cameFrom := "get order status"

	path, err := CheckCourierInProgress(w, r, &cameFrom)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if path != nil {
		http.Redirect(w, r, *path, http.StatusSeeOther)
		return
	}

	orderId, err := strconv.Atoi(r.FormValue("order_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status, err := storage.GetOrderStatus(orderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"status": status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func TakeOrderFromShop(w http.ResponseWriter, r *http.Request) {
	cameFrom := "take order from shop"

	path, err := CheckCourierInProgress(w, r, &cameFrom)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if path != nil {
		http.Redirect(w, r, *path, http.StatusSeeOther)
		return
	}

	orderId, err := strconv.Atoi(r.FormValue("order_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status, err := storage.GetOrderStatus(orderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if status != "order collected" {
		tmpl, err := template.ParseFiles("front/pages/courier/not_yet.html")
		if err != nil {
			http.Error(w, err.Error(), 400)
			log.Fatalf("StartPage: %s", err.Error())
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), 500)
			log.Fatalf("StartPage: %s", err.Error())
		}
		return
	} else {
		err = storage.UpdateOrderStatus(orderId, "order taken from shop")
		if err != nil {
			http.Error(w, err.Error(), 500)
			log.Fatalf("StartPage: %s", err.Error())
		}

		ptrOrderMessage, err := storage.GetFullOrderInfo(orderId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var tookOrderMessage models.OrderMessage
		if ptrOrderMessage != nil {
			tookOrderMessage = *ptrOrderMessage
		}

		err = mq.ProduceMessage(tookOrderMessage, "Order taken from shop")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func FinishDelivery(w http.ResponseWriter, r *http.Request) {
	cameFrom := "finish delivery"

	path, err := CheckCourierInProgress(w, r, &cameFrom)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if path != nil {
		http.Redirect(w, r, *path, http.StatusSeeOther)
		return
	}

	orderId, err := strconv.Atoi(r.FormValue("order_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status, err := storage.GetOrderStatus(orderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if status != "order taken from shop" {
		tmpl, err := template.ParseFiles("front/pages/courier/not_yet.html")
		if err != nil {
			http.Error(w, err.Error(), 400)
			log.Fatalf("StartPage: %s", err.Error())
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), 500)
			log.Fatalf("StartPage: %s", err.Error())
		}
		return
	} else {
		err = storage.UpdateOrderStatus(orderId, "delivered")
		if err != nil {
			http.Error(w, err.Error(), 500)
			log.Fatalf("StartPage: %s", err.Error())
		}

		ptrOrderMessage, err := storage.GetFullOrderInfo(orderId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var deliveredOrderMessage models.OrderMessage
		if ptrOrderMessage != nil {
			deliveredOrderMessage = *ptrOrderMessage
		}

		courierId := GetCourierId(w, r)
		ptrCourier, err := storage.GetState(courierId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var courier models.Courier
		if ptrCourier != nil {
			courier = *ptrCourier
		}

		deliveredAt := time.Now()
		deliveredOrderMessage.DeliveredAt = deliveredAt
		deliveredOrderMessage.Courier = courier

		err = mq.ProduceMessage(deliveredOrderMessage, "Delivered")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = storage.FinishDelivery(courierId, orderId, deliveredAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "http://localhost:8083/courier/delivery_finished", http.StatusSeeOther)
	}
}

func Declined(w http.ResponseWriter, r *http.Request) {
	courierId := GetCourierId(w, r)
	err := storage.DisableInProgress(courierId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func Decline(w http.ResponseWriter, r *http.Request) {
	orderId, err := strconv.Atoi(r.FormValue("order_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = storage.DeclineOrder(orderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	declinedByCourierMessage, err := storage.GetFullOrderInfo(orderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	err = mq.ProduceMessage(*declinedByCourierMessage, "Declined by courier")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("front/pages/courier/decline.html")
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Fatalf("StartPage: %s", err.Error())
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalf("StartPage: %s", err.Error())
	}

	courierId := GetCourierId(w, r)
	err = storage.DisableInProgress(courierId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = storage.DeclinedByCourier(courierId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
