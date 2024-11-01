package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized 1", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte("secret_key"), nil
		})
		if err != nil {
			http.Error(w, "Unauthorized 3", http.StatusUnauthorized)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			http.Error(w, "Unauthorized 4", http.StatusUnauthorized)
			return
		}

		if time.Now().Unix() > int64(claim["exp"].(float64)) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
