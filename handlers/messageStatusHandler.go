package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/types"
)

type MessageStatusHandler interface {
	FindMessageStatusByMessageId(c *gin.Context)
	FindMessageStatusByMessageIds(c *gin.Context)
}

type messageStatusHandler struct {
	messageStatusService services.MessageStatusService
}

func NewMessageStatusHandler(
	messageStatusService services.MessageStatusService,
) *messageStatusHandler {
	return &messageStatusHandler{
		messageStatusService: messageStatusService,
	}
}

func (h *messageStatusHandler) FindMessageStatusByMessageId(c *gin.Context) {
	messageId := c.Param("messageId")
	messageStatus, err := h.messageStatusService.FindMessageStatusByMessageId(messageId)
	if err != nil {
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
		Data:    messageStatus,
	})

}

func (h *messageStatusHandler) FindMessageStatusByMessageIds(c *gin.Context) {
	var messageIds []string
	err := c.BindJSON(&messageIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	messageStatuses, err := h.messageStatusService.FindMessageStatusByMessageIds(messageIds)
	if err != nil {
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
		Data:    messageStatuses,
	})
}
