package controllers

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/saconstant"
	"github.com/tieubaoca/go-chat-server/services"
)

// func FindUserByUsername(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	username, ok := vars["username"]
// 	log.Println(username)
// 	if !ok {
// 		w.WriteHeader(http.StatusBadRequest)
// 		response.Res(w, saconstant.StatusError, nil, "Username is empty")
// 		return
// 	}
// 	user, err := services.FindUserByUsername(username)
// 	if err != nil {
// 		w.WriteHeader(http.StatusNoContent)
// 		response.Res(w, saconstant.StatusError, nil, err.Error())
// 		return
// 	}
// 	response.Res(w, saconstant.StatusSuccess, user, "Find user by username successfully")
// }

// func FindUserById(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, ok := vars["id"]
// 	if !ok {
// 		w.WriteHeader(http.StatusBadRequest)
// 		response.Res(w, saconstant.StatusError, nil, "Id is empty")
// 		return
// 	}
// 	user, err := services.FindUserById(id)
// 	if err != nil {
// 		w.WriteHeader(http.StatusNoContent)
// 		response.Res(w, saconstant.StatusError, nil, err.Error())
// 		return
// 	}
// 	response.Res(w, saconstant.StatusSuccess, user, "Find user by id successfully")
// }

func FindOnlineUsers(w http.ResponseWriter, r *http.Request) {
	users, err := services.FindOnlineUsers()
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	response.Res(w, saconstant.StatusSuccess, users, "Find online users successfully")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	accessTokenString, err := r.Cookie("access-token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	accessToken, err := services.Parse(accessTokenString.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	sessions, err := store.Get(r, "session")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	http.SetCookie(
		w,
		&http.Cookie{
			Name: "access-token",
			Path: "/",
		},
	)
	http.SetCookie(
		w,
		&http.Cookie{
			Name: "refresh-token",
			Path: "/",
		},
	)
	sessions.Values = make(map[interface{}]interface{})
	sessions.Save(r, w)
	services.Logout(accessToken.Claims.(jwt.MapClaims)["preferred_username"].(string), sessions.ID)
	response.Res(w, saconstant.StatusSuccess, nil, "Logout successfully")
}
