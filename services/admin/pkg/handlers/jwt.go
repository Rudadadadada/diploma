package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("randomString")


func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("admin")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "http://localhost:8082/authorization/admin", http.StatusSeeOther)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		tokenStr := c.Value
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Redirect(w, r, "http://localhost:8082/authorization/admin", http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, "http://localhost:8082/authorization/admin", http.StatusSeeOther)
			return
		}

		if !token.Valid {
			http.Redirect(w, r, "http://localhost:8082/authorization/admin", http.StatusSeeOther)
			return
		}
        
        claims, ok := token.Claims.(jwt.MapClaims)
        
        if !ok {
            http.Redirect(w, r, "http://localhost:8082/authorization/admin", http.StatusSeeOther)
			return
        }
        
        scope := claims["scope"].(string)

        if scope != "admin" {
            http.Redirect(w, r, "http://localhost:8082/authorization/admin", http.StatusSeeOther)
			return
        }

		next.ServeHTTP(w, r)
	})
}
