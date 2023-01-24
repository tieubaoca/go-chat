package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tieubaoca/go-chat-server/utils/log"

	"github.com/tieubaoca/go-chat-server/dto/request"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func Authentication(w http.ResponseWriter, r *http.Request) {
	tokenString := utils.GetAccessTokenByReq(r)
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, types.StatusError, nil, types.ErrorTokenEmpty)
		return
	}
	token, err := utils.ParseUnverified(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	response.Res(w, types.StatusSuccess, token.Claims, "")
}

func GetAccessToken(w http.ResponseWriter, r *http.Request) {
	var getAccTokenReq request.GetAccessTokenReq
	err := json.NewDecoder(r.Body).Decode(&getAccTokenReq)
	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	if getAccTokenReq.Username == "" || getAccTokenReq.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, types.StatusError, nil, types.ErrorInvalidInput)
		return
	}
	accessToken, refreshToken, err := services.GetAccessToken(getAccTokenReq.Username, getAccTokenReq.Password)
	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "access-token",
		Value:   accessToken,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "refresh-token",
		Value:   refreshToken,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	})
	_, err = utils.Parse(accessToken)
	if err != nil {
		log.ErrorLogger.Println(err)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}

	response.Res(
		w, types.StatusSuccess,
		bson.M{
			"access-token":  accessToken,
			"refresh-token": refreshToken,
		},
		"",
	)

}

// func RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
// 	accessTokenCookie, err := r.Cookie("access-token")
// 	if err != nil {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		response.Res(w, types.StatusError, nil, err.Error())
// 		return
// 	}
// 	accessToken := accessTokenCookie.Value
// 	if accessToken == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		response.Res(w, types.StatusError, nil, "Access token is empty")
// 		return
// 	}
// 	refreshedAccessToken, err := services.RefreshAccessToken(accessToken)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		response.Res(w, types.StatusError, nil, err.Error())
// 		return
// 	}
// 	http.SetCookie(w, &http.Cookie{
// 		Name:    "access_token",
// 		Value:   refreshedAccessToken,
// 		Expires: time.Now().Add(24 * time.Hour),
// 	})
// 	response.Res(w, types.StatusSuccess, nil, "Refresh access token successfully")
// }
