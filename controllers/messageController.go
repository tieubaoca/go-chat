package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/tieubaoca/go-chat-server/dto/request"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/models"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"
	"github.com/tieubaoca/go-chat-server/utils/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindMessagesByChatRoomId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatRoomId, ok := vars["chatRoomId"]
	if !ok {
		log.ErrorLogger.Println(types.ErrorInvalidInput)
		w.WriteHeader(http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
	}
	chatRoom, err := services.FindChatroomById(chatRoomId)
	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
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

	if !utils.ContainsString(chatRoom.Members, utils.GetSaIdFromToken(token)) {
		log.ErrorLogger.Println(types.ErrorNotRoomMember)
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, types.StatusError, nil, types.ErrorNotRoomMember)
		return
	}

	messages, err := services.FindMessagesByChatRoomId(chatRoomId)
	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	response.Res(w, types.StatusSuccess, messages, "")
}

func PaginationMessagesByChatRoomId(w http.ResponseWriter, r *http.Request) {
	var pagination request.MessagePaginationReq
	err := json.NewDecoder(r.Body).Decode(&pagination)
	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, types.StatusError, nil, err.Error())
	}

	chatRoom, err := services.FindChatroomById(pagination.ChatRoomId)
	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
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

	if !utils.ContainsString(chatRoom.Members, utils.GetSaIdFromToken(token)) {
		log.ErrorLogger.Println(types.ErrorNotRoomMember)
		w.WriteHeader(http.StatusUnauthorized)
		response.Res(w, types.StatusError, nil, types.ErrorNotRoomMember)
		return
	}

	messages, err := services.PaginationMessagesByChatRoomId(
		pagination.ChatRoomId,
		pagination.Limit,
		pagination.Skip,
	)
	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	response.Res(w, types.StatusSuccess, messages, "")
}

func InsertMessage(w http.ResponseWriter, r *http.Request) {
	var message models.Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	result, err := services.InsertMessage(bson.M{
		"chatroom": message.Chatroom,
		"sender":   message.Sender,
		"content":  message.Content,
		"createAt": primitive.NewDateTimeFromTime(time.Now()),
	})
	if err != nil {
		log.ErrorLogger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		response.Res(w, types.StatusError, nil, err.Error())
		return
	}
	response.Res(w, types.StatusSuccess, result, "")
}
