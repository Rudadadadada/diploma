package handlers

import (
	"diploma/services/authorization/pkg/models"
	"diploma/services/authorization/pkg/storage"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func CustomerRegistration(w http.ResponseWriter, r *http.Request) {
	user := models.Customer{
		Name:     r.FormValue("first-name"),
		Surname:  r.FormValue("last-name"),
		Email:    r.FormValue("email"),
		HashPassword: "",
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.HashPassword = string(hashPassword)

	err = storage.CustomerRegistration(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("front/pages/authorization/success.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := SuccessData{
		Title: "Вы успешно зарегистрировались. Пройдите авторизацию",
		Path:  "/authorization/customer",
	}

	if err = tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func CourierRegistration(w http.ResponseWriter, r *http.Request) {
	user := models.Courier{
		Name:     r.FormValue("first-name"),
		Surname:  r.FormValue("last-name"),
		Email:    r.FormValue("email"),
		HashPassword: "",
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.HashPassword = string(hashPassword)

	err = storage.CourierRegistration(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("front/pages/authorization/success.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := SuccessData{
		Title: "Вы успешно зарегистрировались. Пройдите авторизацию",
		Path:  "/authorization/courier",
	}

	if err = tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}