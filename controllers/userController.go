package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/tieubaoca/go-chat-server/utils/log"

	"github.com/tieubaoca/go-chat-server/dto/request"
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

func Logout(w http.ResponseWriter, r *http.Request) {

	token, err := utils.ParseUnverified(utils.GetAccessTokenByReq(r))
	if err != nil {
		log.ErrorLogger.Println(err)
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
	response.Res(w, types.StatusSuccess, nil, "")
}

func PaginationOnlineFriend(w http.ResponseWriter, r *http.Request) {
	tokenString := utils.GetAccessTokenByReq(r)
	token, err := utils.ParseUnverified(tokenString)
	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	saId := utils.GetSaIdFromToken(token)

	var paginationReq request.PaginationOnlineFriendReq
	err = json.NewDecoder(r.Body).Decode(&paginationReq)
	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}

	users, err := services.PaginationOnlineFriends(saId, paginationReq)
	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusNoContent)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	response.Res(w, types.StatusSuccess, users, "")
}
