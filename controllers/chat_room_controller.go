package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tieubaoca/go-chat-server/dto"
	"github.com/tieubaoca/go-chat-server/saconstant"
	"github.com/tieubaoca/go-chat-server/services"
)

func FindChatRoomById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		dto.Res(w, saconstant.StatusError, nil, "Id is empty")
		return
	}
	chatRoom, err := services.FindChatroomById(id)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		dto.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	dto.Res(w, saconstant.StatusSuccess, chatRoom, "Find chat room by id successfully")
}

func FindChatRoomsByMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	member, ok := vars["member"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		dto.Res(w, saconstant.StatusError, nil, "Member is empty")
		return
	}
	chatRooms, err := services.FindChatroomsByMember(member)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		dto.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	dto.Res(w, saconstant.StatusSuccess, chatRooms, "Find chat rooms by member successfully")
}

func FindDMByMembers(w http.ResponseWriter, r *http.Request) {
	var members []string
	err := json.NewDecoder(r.Body).Decode(&members)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		dto.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	if len(members) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		dto.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}

	chatRoom, err := services.FindDMByMembers(members)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		dto.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	dto.Res(w, saconstant.StatusSuccess, chatRoom, "Find chat room by members successfully")

}

func InsertChatRoom(w http.ResponseWriter, r *http.Request) {
	var chatRoom interface{}
	err := json.NewDecoder(r.Body).Decode(&chatRoom)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		dto.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	result, err := services.InsertChatroom(chatRoom)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		dto.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	dto.Res(w, saconstant.StatusSuccess, result, "Insert chat room successfully")
}
