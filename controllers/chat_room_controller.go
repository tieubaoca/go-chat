package controllers

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/models"
	"github.com/tieubaoca/go-chat-server/saconstant"
	"github.com/tieubaoca/go-chat-server/services"
	"go.mongodb.org/mongo-driver/bson"
)

func FindChatRoomById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, saconstant.StatusError, nil, "Id is empty")
		return
	}
	chatRoom, err := services.FindChatroomById(id)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	response.Res(w, saconstant.StatusSuccess, chatRoom, "Find chat room by id successfully")
}

func FindChatRoomsByMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	member, ok := vars["member"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, saconstant.StatusError, nil, "Member is empty")
		return
	}
	chatRooms, err := services.FindChatroomsByMember(member)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	response.Res(w, saconstant.StatusSuccess, chatRooms, "Find chat rooms by member successfully")
}

func FindDMByMembers(w http.ResponseWriter, r *http.Request) {
	var member string
	err := json.NewDecoder(r.Body).Decode(&member)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	if member == "" {
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	session, err := store.Get(r, "session")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}

	members := []string{member, session.Values["username"].(string)}

	chatRoom, err := services.FindDMByMembers(members)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	response.Res(w, saconstant.StatusSuccess, chatRoom, "Find chat room by members successfully")

}

func FindGroupsByMembers(w http.ResponseWriter, r *http.Request) {
	var members []string
	err := json.NewDecoder(r.Body).Decode(&members)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	if len(members) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}

}

func CreateNewGroupChat(w http.ResponseWriter, r *http.Request) {
	var chatRoom models.Chatroom
	err := json.NewDecoder(r.Body).Decode(&chatRoom)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	if chatRoom.Type != models.ChatroomTypeGroup {
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, saconstant.StatusError, nil, "Chat room type is not group")
		return
	}
	if chatRoom.Owner == "" {
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, saconstant.StatusError, nil, "Chat room owner is empty")
		return
	}

	result, err := services.InsertChatroom(
		bson.M{
			"name":    chatRoom.Name,
			"type":    chatRoom.Type,
			"owner":   chatRoom.Owner,
			"members": chatRoom.Members,
		},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	response.Res(w, saconstant.StatusSuccess, result, "Insert chat room successfully")
}

func CreateDMRoom(w http.ResponseWriter, r *http.Request) {
	var member string
	err := json.NewDecoder(r.Body).Decode(&member)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	if member == "" {
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	session, err := store.Get(r, "session")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	members := []string{
		member,
		session.Values["username"].(string),
	}
	sort.Strings(members)
	result, err := services.InsertChatroom(
		bson.M{
			"name":    members[0] + "-" + members[1],
			"type":    models.ChatroomTypeDM,
			"members": members,
		},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	response.Res(w, saconstant.StatusSuccess, result, "Insert chat room successfully")
}
