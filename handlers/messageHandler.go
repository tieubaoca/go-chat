package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tieubaoca/go-chat-server/dto/request"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"
	"github.com/tieubaoca/go-chat-server/utils/log"
)

type MessageHandler interface {
	FindMessagesByChatRoomId(c *gin.Context)
	PaginationMessagesByChatRoomId(c *gin.Context)
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
	_, err = h.chatRoomService.FindChatRoomById(saId, chatRoomId)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}

	messages, err := h.messageService.FindMessagesByChatRoomId(saId, chatRoomId)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, response.ResponseData{
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
	_, err = h.chatRoomService.FindChatRoomById(saId, pagination.ChatRoomId)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}

	messages, err := h.messageService.PaginationMessagesByChatRoomId(
		saId,
		pagination,
	)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, response.ResponseData{
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
