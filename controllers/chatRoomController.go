package controllers

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/models"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"
	"github.com/tieubaoca/go-chat-server/utils/log"
)

func FindChatRoomById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.ErrorLogger.Println(types.ErrorInvalidInput)
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, types.StatusError, nil, types.ErrorInvalidInput)
		return
	}
	chatRoom, err := services.FindChatRoomById(id)
	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusNoContent)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	token, _ := utils.ParseUnverified(utils.GetAccessTokenByReq(r))
	if !utils.ContainsString(chatRoom.Members, utils.GetSaIdFromToken(token)) {
		log.ErrorLogger.Println(types.ErrorNotRoomMember)
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, types.StatusError, nil, types.ErrorNotRoomMember)
		return
	}
	response.Res(w, types.StatusSuccess, chatRoom, "")
}

func FindChatRooms(w http.ResponseWriter, r *http.Request) {
	token, err := utils.ParseUnverified(utils.GetAccessTokenByReq(r))
	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	chatRooms, err := services.FindChatRoomsByMember(utils.GetSaIdFromToken(token))
	if err != nil {
		log.ErrorLogger.Println(err)
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
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	if member == "" {
		log.ErrorLogger.Println(types.ErrorInvalidInput)
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, types.StatusError, nil, types.ErrorInvalidInput)
		return
	}
	token, err := utils.ParseUnverified(utils.GetAccessTokenByReq(r))
	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}

	members := []string{member, utils.GetSaIdFromToken(token)}

	chatRoom, err := services.FindDMByMembers(members)
	if err != nil {
		log.ErrorLogger.Println(err)
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
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}

	chatRooms, err := services.FindGroupsByMembers(append(members, utils.GetSaIdFromToken(token)))
	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	response.Res(w, types.StatusSuccess, chatRooms, "")
}

func CreateNewGroupChat(w http.ResponseWriter, r *http.Request) {

	var chatRoom models.ChatRoom
	err := json.NewDecoder(r.Body).Decode(&chatRoom)

	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	chatRoom.Type = models.ChatRoomTypeGroup
	token, _ := utils.ParseUnverified(utils.GetAccessTokenByReq(r))
	chatRoom.Owner = utils.GetSaIdFromToken(token)

	result, err := services.InsertChatRoom(chatRoom)
	if err != nil {
		log.ErrorLogger.Println(err)
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
		log.ErrorLogger.Println(err)
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
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}

	members := []string{
		member,
		utils.GetSaIdFromToken(token),
	}
	sort.Strings(members)
	result, err := services.InsertChatRoom(
		models.ChatRoom{
			Name:    members[0] + "-" + members[1],
			Type:    models.ChatRoomTypeDM,
			Members: members,
		},
	)
	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	response.Res(w, types.StatusSuccess, result, "")
}
