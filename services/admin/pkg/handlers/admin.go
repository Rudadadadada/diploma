package handlers

import (
	"diploma/services/admin/pkg/redis"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("admin")
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	tokenStr := c.Value
	redis.RedisClient.Del(tokenStr)

	http.SetCookie(w, &http.Cookie{
		Name:     "admin",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1,
	})

	http.Redirect(w, r, "http://localhost:8082/authorization/admin", http.StatusSeeOther)
}
