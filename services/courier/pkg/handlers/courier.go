package handlers

import (
	"diploma/services/courier/pkg/redis"
	"diploma/services/courier/pkg/storage"
	"encoding/json"
	"net/http"

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

	err := storage.AddCourier(courierId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = storage.SetActive(courierId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "http://localhost:8083/courier", http.StatusSeeOther)
}

func GetState(w http.ResponseWriter, r *http.Request) {
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

	response := map[string]bool{"active": active}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Ошибка при форматировании ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}