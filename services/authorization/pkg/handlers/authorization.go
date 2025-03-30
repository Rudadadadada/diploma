package handlers

import (
	"diploma/services/authorization/pkg/models"
	"diploma/services/authorization/pkg/redis"
	"diploma/services/authorization/pkg/storage"
	"html/template"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("randomString")


func CustomerAuthorization(w http.ResponseWriter, r *http.Request) {
	customer := models.Customer{
		Email:        r.FormValue("customer"),
		HashPassword: r.FormValue("password"),
	}

	dbCustomer, err := storage.CustomerAuthorization(customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbCustomer.HashPassword), []byte(customer.HashPassword)); err != nil {
        t, _ := template.ParseFiles("front/pages/authorization/customer_authorization.html")
        t.Execute(w, map[string]string{"error": "Неверная почта или пароль"})
        return
    }

	token := jwt.New(jwt.SigningMethodHS256)
	claimes := token.Claims.(jwt.MapClaims)
	claimes["scope"] = "customer"
	claimes["customer_id"] = dbCustomer.Id

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	ttl := time.Minute * 15
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
	courier := models.Courier{
		Email:        r.FormValue("courier"),
		HashPassword: r.FormValue("password"),
	}

	dbCourier, err := storage.CourierAuthorization(courier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbCourier.HashPassword), []byte(courier.HashPassword)); err != nil {
        t, _ := template.ParseFiles("front/pages/authorization/courier_authorization.html")
        t.Execute(w, map[string]string{"error": "Неверная почта или пароль"})
        return
    }

	token := jwt.New(jwt.SigningMethodHS256)
	claimes := token.Claims.(jwt.MapClaims)
	claimes["scope"] = "courier"
	claimes["courier_id"] = dbCourier.Id

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	ttl := time.Minute * 15
	err = redis.SetKeyWithTTL(tokenString, tokenString, ttl)
    if err != nil {
        http.Error(w, "Could not save token to Redis", http.StatusInternalServerError)
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

	if err := bcrypt.CompareHashAndPassword([]byte(dbHashPassword), []byte(admin.HashPassword)); err != nil {
        t, _ := template.ParseFiles("front/pages/authorization/admin_authorization.html")
        t.Execute(w, map[string]string{"error": "Неверная почта или пароль"})
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

	ttl := time.Minute * 15
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
