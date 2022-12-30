package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tieubaoca/go-chat-server/dto"
	"github.com/tieubaoca/go-chat-server/saconstant"
	"github.com/tieubaoca/go-chat-server/services"
)

func FindUserByUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, ok := vars["username"]
	log.Println(username)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		dto.Res(w, saconstant.StatusError, nil, "Username is empty")
		return
	}
	user, err := services.FindUserByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		dto.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	dto.Res(w, saconstant.StatusSuccess, user, "Find user by username successfully")
}

func FindUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		dto.Res(w, saconstant.StatusError, nil, "Id is empty")
		return
	}
	user, err := services.FindUserById(id)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		dto.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	dto.Res(w, saconstant.StatusSuccess, user, "Find user by id successfully")
}

func FindOnlineUsers(w http.ResponseWriter, r *http.Request) {
	users, err := services.FindOnlineUsers()
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		dto.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	dto.Res(w, saconstant.StatusSuccess, users, "Find online users successfully")
}
