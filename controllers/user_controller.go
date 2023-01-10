package controllers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"
)

// func FindUserByUsername(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	username, ok := vars["username"]
// 	log.Println(username)
// 	if !ok {
// 		w.WriteHeader(http.StatusBadRequest)
// 		response.Res(w, types.StatusError, nil, "Username is empty")
// 		return
// 	}
// 	user, err := services.FindUserByUsername(username)
// 	if err != nil {
// 		w.WriteHeader(http.StatusNoContent)
// 		response.Res(w, types.StatusError, nil, err.Error())
// 		return
// 	}
// 	response.Res(w, types.StatusSuccess, user, "Find user by username successfully")
// }

// func FindUserById(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, ok := vars["id"]
// 	if !ok {
// 		w.WriteHeader(http.StatusBadRequest)
// 		response.Res(w, types.StatusError, nil, "Id is empty")
// 		return
// 	}
// 	user, err := services.FindUserById(id)
// 	if err != nil {
// 		w.WriteHeader(http.StatusNoContent)
// 		response.Res(w, types.StatusError, nil, err.Error())
// 		return
// 	}
// 	response.Res(w, types.StatusSuccess, user, "Find user by id successfully")
// }

func FindOnlineFriends(w http.ResponseWriter, r *http.Request) {
	tokenString := utils.GetAccessTokenByReq(r)
	token, err := utils.ParseUnverified(tokenString)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	saId := utils.GetSaIdFromToken(token)

	users, err := services.FindOnlineFriends(saId)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusNoContent)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	response.Res(w, types.StatusSuccess, users, "Find online users successfully")
}

func Logout(w http.ResponseWriter, r *http.Request) {

	token, err := utils.ParseUnverified(utils.GetAccessTokenByReq(r))
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, types.StatusError, nil, err.Error())
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
	services.Logout(utils.GetSaIdFromToken(token), utils.GetSessionIdFromToken(token))
	response.Res(w, types.StatusSuccess, nil, "Logout successfully")
}
