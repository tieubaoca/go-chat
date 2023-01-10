package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/tieubaoca/go-chat-server/utils"
)

func JwtMiddleware(next http.HandlerFunc, requireRole string) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		accessTokenString := utils.GetAccessTokenByReq(r)
		if accessTokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token, err := utils.Parse(accessTokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		account := token.Claims.(jwt.MapClaims)["resource_access"].(map[string]interface{})["account"].(map[string]interface{})
		for _, role := range account["roles"].([]interface{}) {
			if role == requireRole {
				next(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusMethodNotAllowed)

	})
}
