package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/tieubaoca/go-chat-server/dto/request"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/models"
	"github.com/tieubaoca/go-chat-server/saconstant"
	"github.com/tieubaoca/go-chat-server/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindMessagesByChatRoomId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatRoomId, ok := vars["chatRoomId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
	}
	messages, err := services.FindMessagesByChatRoomId(chatRoomId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	response.Res(w, saconstant.StatusSuccess, messages, "Get messages successfully")
}

func PaginationMessagesByChatRoomId(w http.ResponseWriter, r *http.Request) {
	var pagination request.MessagePaginationReq
	err := json.NewDecoder(r.Body).Decode(&pagination)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, saconstant.StatusError, nil, err.Error())
	}
	messages, err := services.PaginationMessagesByChatRoomId(
		pagination.ChatRoomId,
		pagination.Limit,
		pagination.Skip,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	response.Res(w, saconstant.StatusSuccess, messages, "Get messages successfully")
}

func InsertMessage(w http.ResponseWriter, r *http.Request) {
	var message models.Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	result, err := services.InsertMessage(bson.M{
		"chatroom": message.Chatroom,
		"sender":   message.Sender,
		"content":  message.Content,
		"createAt": primitive.DateTime(time.Now().Unix() * 1000),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Res(w, saconstant.StatusError, nil, err.Error())
		return
	}
	response.Res(w, saconstant.StatusSuccess, result, "Insert message successfully")
}
