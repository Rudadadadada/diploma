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
	courierId := GetCourierId(w, r)

	err := storage.SetActive(courierId)

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

	var active bool
	if ptr != nil {
		active = *ptr
	}

	response := map[string]bool{"active": active}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Ошибка при форматировании ответа", http.StatusInternalServerError)
		return
	}

	courierState := models.CourierState{
		CourierId: courierId,
		State:     active,
	}
	err = mq.ProduceState(courierState, "Courier state")
	if err != nil {
		http.Error(w, "Ошибка при форматировании ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func CheckCourierState(w http.ResponseWriter, r *http.Request, redirectPath string) {
	courierId := GetCourierId(w, r)

	ptr, err := storage.GetState(courierId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var active bool
	if ptr != nil {
		active = *ptr
	}

	if !active {
		http.Redirect(w, r, redirectPath, http.StatusSeeOther)
		return
	}
}

func TakeOrder(w http.ResponseWriter, r *http.Request) {
	CheckCourierState(w, r, "http://localhost:8083/courier")

	orderId, err := strconv.Atoi(r.FormValue("order_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ptr, err := storage.CheckOrderTaken(orderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	ptrOrderMessage, err := storage.GetFullOrderInfo(orderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var tookOrderMessage models.OrderMessage
	if ptr != nil {
		tookOrderMessage = *ptrOrderMessage
	}

	tookOrderMessage.Status = "preparing"
	tookOrderMessage.DeliveryStartedAt = time.Now()
	tookOrderMessage.Courier = models.Courier{Id: uint(courierId)}
	
	err = mq.ProduceMessage(tookOrderMessage, "Order taken")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("front/pages/courier/order_taken.html")
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
