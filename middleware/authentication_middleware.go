package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/tieubaoca/go-chat-server/services"
)

func JwtMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		accessTokenCookie, err := r.Cookie("access-token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if accessTokenCookie.Value == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := services.Parse(accessTokenCookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		r.Header.Set("username", token.Claims.(jwt.MapClaims)["preferred_username"].(string))
		next(w, r)
	})
}
