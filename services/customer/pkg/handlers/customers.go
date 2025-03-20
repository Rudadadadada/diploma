package handlers

import (
	"diploma/services/customer/pkg/redis"
	"diploma/services/customer/pkg/storage"
	"html/template"
	"net/http"
	"strconv"
)

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
		MaxAge: -1,
	})

	http.Redirect(w, r, "http://localhost:8082/authorization/customer", http.StatusSeeOther)
}

func SelectProductsByCategoryId(w http.ResponseWriter, r *http.Request) {
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

	tmpl, err := template.ParseFiles("front/pages/customer/select_products.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err = tmpl.Execute(w, products); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}		
}
