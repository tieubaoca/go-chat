package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/tieubaoca/go-chat-server/dto/request"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/saconstant"
	"github.com/tieubaoca/go-chat-server/services"
)

func Authentication(w http.ResponseWriter, r *http.Request) {
	sessions, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Res(w, saconstant.StatusSuccess, sessions.Values["username"], "Authentication successfully")
}

func GetAccessToken(w http.ResponseWriter, r *http.Request) {
	var getAccTokenReq request.GetAccessTokenReq
	err := json.NewDecoder(r.Body).Decode(&getAccTokenReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	if getAccTokenReq.Username == "" || getAccTokenReq.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, saconstant.StatusError, nil, "Username or password is empty")
		return
	}
	accessToken, refreshToken, err := services.GetAccessToken(getAccTokenReq.Username, getAccTokenReq.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Res(w, saconstant.StatusError, nil, err.Error())
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
	session, err := store.Get(r, "session")
	if err != nil {
		response.Res(w, saconstant.StatusError, nil, "Get session failed")
		return
	}
	token, err := services.Parse(accessToken)
	if err != nil {
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}

	session.Values["username"] = token.Claims.(jwt.MapClaims)["preferred_username"].(string)

	// Save it before we write to the response/return from the handler.
	err = session.Save(r, w)
	if err != nil {
		response.Res(w, saconstant.StatusError, nil, "Save session failed")
		return
	}
	response.Res(w, saconstant.StatusSuccess, nil, "Get access token successfully")

}

// func RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
// 	accessTokenCookie, err := r.Cookie("access-token")
// 	if err != nil {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		response.Res(w, saconstant.StatusError, nil, err.Error())
// 		return
// 	}
// 	accessToken := accessTokenCookie.Value
// 	if accessToken == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		response.Res(w, saconstant.StatusError, nil, "Access token is empty")
// 		return
// 	}
// 	refreshedAccessToken, err := services.RefreshAccessToken(accessToken)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		response.Res(w, saconstant.StatusError, nil, err.Error())
// 		return
// 	}
// 	http.SetCookie(w, &http.Cookie{
// 		Name:    "access_token",
// 		Value:   refreshedAccessToken,
// 		Expires: time.Now().Add(24 * time.Hour),
// 	})
// 	response.Res(w, saconstant.StatusSuccess, nil, "Refresh access token successfully")
// }
