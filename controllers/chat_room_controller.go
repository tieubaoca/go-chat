package controllers

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/models"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"
)

func FindChatRoomById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.Error(types.ErrorInvalidInput)
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, types.StatusError, nil, types.ErrorInvalidInput)
		return
	}
	chatRoom, err := services.FindChatroomById(id)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusNoContent)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	response.Res(w, types.StatusSuccess, chatRoom, "")
}

func FindChatRooms(w http.ResponseWriter, r *http.Request) {
	token, err := utils.ParseUnverified(utils.GetAccessTokenByReq(r))
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	chatRooms, err := services.FindChatroomsByMember(utils.GetUsernameFromToken(token))
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusNoContent)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	response.Res(w, types.StatusSuccess, chatRooms, "")
}

func FindDMByMembers(w http.ResponseWriter, r *http.Request) {
	var member string
	err := json.NewDecoder(r.Body).Decode(&member)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	if member == "" {
		log.Error(types.ErrorInvalidInput)
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, types.StatusError, nil, types.ErrorInvalidInput)
		return
	}
	token, err := utils.ParseUnverified(utils.GetAccessTokenByReq(r))
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}

	members := []string{member, utils.GetUsernameFromToken(token)}

	chatRoom, err := services.FindDMByMembers(members)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	response.Res(w, types.StatusSuccess, chatRoom, "")

}

func FindGroupsByMembers(w http.ResponseWriter, r *http.Request) {
	var members []string
	err := json.NewDecoder(r.Body).Decode(&members)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		members = make([]string, 0)
	}
	token, err := utils.ParseUnverified(utils.GetAccessTokenByReq(r))
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}

	chatrooms, err := services.FindGroupsByMembers(append(members, utils.GetUsernameFromToken(token)))
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	response.Res(w, types.StatusSuccess, chatrooms, "")
}

func CreateNewGroupChat(w http.ResponseWriter, r *http.Request) {
	var chatRoom models.Chatroom
	err := json.NewDecoder(r.Body).Decode(&chatRoom)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	if chatRoom.Type != models.ChatroomTypeGroup {
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, types.StatusError, nil, types.ErrorInvalidInput)
		return
	}
	if chatRoom.Owner == "" {
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, types.StatusError, nil, types.ErrorInvalidInput)
		return
	}

	result, err := services.InsertChatroom(chatRoom)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	response.Res(w, types.StatusSuccess, result, "")
}

func CreateDMRoom(w http.ResponseWriter, r *http.Request) {
	var member string
	err := json.NewDecoder(r.Body).Decode(&member)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	if member == "" {
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	token, err := utils.ParseUnverified(utils.GetAccessTokenByReq(r))
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}

	members := []string{
		member,
		utils.GetUsernameFromToken(token),
	}
	sort.Strings(members)
	result, err := services.InsertChatroom(
		models.Chatroom{
			Name:    members[0] + "-" + members[1],
			Type:    models.ChatroomTypeDM,
			Members: members,
		},
	)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	response.Res(w, types.StatusSuccess, result, "")
}
