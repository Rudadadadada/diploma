package handlers

import (
	"diploma/services/authorization/pkg/models"
	"diploma/services/authorization/pkg/redis"
	"diploma/services/authorization/pkg/storage"
	_ "log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("randomString")


func CustomerAuthorization(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		Email:        r.FormValue("email"),
		HashPassword: r.FormValue("password"),
	}

	dbUser, err := storage.CustomerAuthorization(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.HashPassword), []byte(user.HashPassword))
	if err != nil {
		http.Error(w, "Неверная почта или пароль", http.StatusBadRequest)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claimes := token.Claims.(jwt.MapClaims)
	claimes["scope"] = "customer"
	claimes["customer_id"] = dbUser.Id

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}


	ttl := time.Minute * 5
	err = redis.SetKeyWithTTL(tokenString, tokenString, ttl)
    if err != nil {
        http.Error(w, "Could not save token to Redis", http.StatusInternalServerError)
        return
    }

	http.SetCookie(w, &http.Cookie{
		Name:     "customer",
		Value:    tokenString,
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "http://localhost:8081/customer", http.StatusSeeOther)
}

func CourierAuthorization(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		Email:        r.FormValue("email"),
		HashPassword: r.FormValue("password"),
	}

	dbUser, err := storage.CustomerAuthorization(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.HashPassword), []byte(user.HashPassword))
	if err != nil {
		http.Error(w, "Неверная почта или пароль", http.StatusBadRequest)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claimes := token.Claims.(jwt.MapClaims)
	claimes["scope"] = "courier"
	claimes["courier_id"] = dbUser.Id

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "courier",
		Value:    tokenString,
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "http://localhost:8083/courier", http.StatusSeeOther)
}

func AdminAuthorization(w http.ResponseWriter, r *http.Request) {
	admin := models.Admin{
		Admin:        r.FormValue("admin"),
		HashPassword: r.FormValue("password"),
	}

	dbHashPassword, err := storage.AdminAuthorization(admin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbHashPassword), []byte(admin.HashPassword))
	if err != nil {
		http.Error(w, "Неверная почта или пароль", http.StatusBadRequest)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claimes := token.Claims.(jwt.MapClaims)
	claimes["scope"] = "admin"
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	ttl := time.Minute * 5
	err = redis.SetKeyWithTTL(tokenString, tokenString, ttl)
    if err != nil {
        http.Error(w, "Could not save token to Redis", http.StatusInternalServerError)
        return
    }

	http.SetCookie(w, &http.Cookie{
		Name:     "admin",
		Value:    tokenString,
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "http://localhost:8080/admin", http.StatusSeeOther)
}
