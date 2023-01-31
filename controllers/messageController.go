package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/tieubaoca/go-chat-server/dto/request"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/models"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"
	"github.com/tieubaoca/go-chat-server/utils/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindMessagesByChatRoomId(c *gin.Context) {
	vars := mux.Vars(c.Request)
	chatRoomId, ok := vars["chatRoomId"]
	if !ok {
		log.ErrorLogger.Println(types.ErrorInvalidInput)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorInvalidInput,
			Data:    "",
		})
	}
	chatRoom, err := services.FindChatRoomById(chatRoomId)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusNoContent, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	token, err := utils.ParseUnverified(utils.GetAccessTokenByReq(c.Request))
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}

	if !utils.ContainsString(chatRoom.Members, utils.GetSaIdFromToken(token)) {
		log.ErrorLogger.Println(types.ErrorNotRoomMember)
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorNotRoomMember,
			Data:    "",
		})
		return
	}

	messages, err := services.FindMessagesByChatRoomId(chatRoomId)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusNoContent, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	c.JSON(http.StatusOK, response.ResponseData{
		Status:  types.StatusSuccess,
		Message: "",
		Data:    messages,
	})
}

func PaginationMessagesByChatRoomId(c *gin.Context) {
	var pagination request.MessagePaginationReq
	err := json.NewDecoder(c.Request.Body).Decode(&pagination)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
	}

	chatRoom, err := services.FindChatRoomById(pagination.ChatRoomId)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusNoContent, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	token, err := utils.ParseUnverified(utils.GetAccessTokenByReq(c.Request))
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}

	if !utils.ContainsString(chatRoom.Members, utils.GetSaIdFromToken(token)) {
		log.ErrorLogger.Println(types.ErrorNotRoomMember)
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorNotRoomMember,
			Data:    "",
		})
		return
	}

	messages, err := services.PaginationMessagesByChatRoomId(
		pagination.ChatRoomId,
		pagination.Limit,
		pagination.Skip,
	)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusNoContent, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	c.JSON(http.StatusOK, response.ResponseData{
		Status:  types.StatusSuccess,
		Message: "",
		Data:    messages,
	})
}

func InsertMessage(c *gin.Context) {
	var message models.Message
	err := json.NewDecoder(c.Request.Body).Decode(&message)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	message.CreateAt = primitive.NewDateTimeFromTime(time.Now())
	result, err := services.InsertMessage(message)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusNoContent, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	c.JSON(http.StatusOK, response.ResponseData{
		Status:  types.StatusSuccess,
		Message: "",
		Data:    result,
	})
}
