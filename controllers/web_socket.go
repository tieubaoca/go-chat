package controllers

import (
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt"
	"github.com/tieubaoca/go-chat-server/dto"
	"github.com/tieubaoca/go-chat-server/saconstant"
	"github.com/tieubaoca/go-chat-server/services"
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	accessTokenCookie, err := r.Cookie("access-token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		dto.Res(w, saconstant.StatusError, nil, err.Error())

		return
	}
	tokenString := accessTokenCookie.Value
	log.Println(tokenString)
	if govalidator.IsNull(tokenString) {
		w.WriteHeader(http.StatusBadRequest)
		dto.Res(w, saconstant.StatusError, nil, err.Error())

		return
	}
	token, err := services.Parse(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		dto.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	services.HandleWebSocket(w, r, token.Claims.(jwt.MapClaims)["preferred_username"].(string))
}
