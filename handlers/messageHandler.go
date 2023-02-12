package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tieubaoca/go-chat-server/dto/request"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/models"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"
	"github.com/tieubaoca/go-chat-server/utils/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageHandler interface {
	FindMessagesByChatRoomId(c *gin.Context)
	PaginationMessagesByChatRoomId(c *gin.Context)
	InsertMessage(c *gin.Context)
}

type messageHandler struct {
	messageService  services.MessageService
	chatRoomService services.ChatRoomService
}

func NewMessageHandler(
	messageService services.MessageService,
	chatRoomService services.ChatRoomService,
) *messageHandler {
	return &messageHandler{
		messageService:  messageService,
		chatRoomService: chatRoomService,
	}
}

func (h *messageHandler) FindMessagesByChatRoomId(c *gin.Context) {
	chatRoomId := c.Param("chatRoomId")
	if chatRoomId == "" {
		log.ErrorLogger.Println(types.ErrorInvalidInput)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorInvalidInput,
			Data:    "",
		})
	}
	chatRoom, err := h.chatRoomService.FindChatRoomById(chatRoomId)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusNoContent, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}

	saId, err := utils.GetSaIdFromToken(utils.GetAccessTokenByReq(c.Request))
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}

	if !utils.ContainsString(chatRoom.Members, saId) {
		log.ErrorLogger.Println(types.ErrorNotRoomMember)
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorNotRoomMember,
			Data:    "",
		})
		return
	}

	messages, err := h.messageService.FindMessagesByChatRoomId(chatRoomId)
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

func (h *messageHandler) PaginationMessagesByChatRoomId(c *gin.Context) {
	var pagination request.MessagePaginationReq
	err := c.ShouldBindJSON(&pagination)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
	}

	chatRoom, err := h.chatRoomService.FindChatRoomById(pagination.ChatRoomId)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusNoContent, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	saId, err := utils.GetSaIdFromToken(utils.GetAccessTokenByReq(c.Request))
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	if !utils.ContainsString(chatRoom.Members, saId) {
		log.ErrorLogger.Println(types.ErrorNotRoomMember)
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorNotRoomMember,
			Data:    "",
		})
		return
	}

	messages, err := h.messageService.PaginationMessagesByChatRoomId(
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

func (h *messageHandler) InsertMessage(c *gin.Context) {
	var message models.Message
	err := c.ShouldBindJSON(&message)
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
	result, err := h.messageService.InsertMessage(message)
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
