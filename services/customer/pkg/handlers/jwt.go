package handlers

import (
	"diploma/services/customer/pkg/redis"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

var jwtKey = []byte("randomString")

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("customer")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "http://localhost:8082/authorization/customer", http.StatusSeeOther)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		tokenStr := c.Value
		exists, err := redis.RedisClient.Exists(tokenStr).Result()

		if err != nil {
			http.Error(w, "Error checking token in Redis", http.StatusInternalServerError)
			return
		}

		if exists == 0 {
			http.SetCookie(w, &http.Cookie{
				Name:     "customer",
				Value:    "",
				HttpOnly: true,
				Path:     "/",
				MaxAge:   -1,
			})

			http.Redirect(w, r, "http://localhost:8082/authorization/customer", http.StatusSeeOther)
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Redirect(w, r, "http://localhost:8082/authorization/customer", http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, "http://localhost:8082/authorization/customer", http.StatusSeeOther)
			return
		}

		if !token.Valid {
			http.Redirect(w, r, "http://localhost:8082/authorization/customer", http.StatusSeeOther)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			http.Redirect(w, r, "http://localhost:8082/authorization/customer", http.StatusSeeOther)
			return
		}

		scope := claims["scope"].(string)

		if scope != "customer" {
			http.Redirect(w, r, "http://localhost:8082/authorization/customer", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
