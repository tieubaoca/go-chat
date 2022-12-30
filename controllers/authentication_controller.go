package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tieubaoca/go-chat-server/dto"
	"github.com/tieubaoca/go-chat-server/saconstant"
	"github.com/tieubaoca/go-chat-server/services"
)

func Authentication(w http.ResponseWriter, r *http.Request) {
	dto.Res(w, saconstant.StatusSuccess, nil, "Authentication successfully")

}

func GetAccessToken(w http.ResponseWriter, r *http.Request) {
	var getAccTokenReq dto.GetAccessTokenReq
	err := json.NewDecoder(r.Body).Decode(&getAccTokenReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		dto.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	if getAccTokenReq.Username == "" || getAccTokenReq.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		dto.Res(w, saconstant.StatusError, nil, "Username or password is empty")
		return
	}
	accessToken, err := services.GetAccessToken(getAccTokenReq.Username, getAccTokenReq.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		dto.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "access-token",
		Value:   accessToken,
		Expires: time.Now().Add(24 * time.Hour),
	})
	dto.Res(w, saconstant.StatusSuccess, nil, "Get access token successfully")

}
