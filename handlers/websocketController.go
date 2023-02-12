package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils/log"

	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/utils"
)

type WebSocketHandler interface {
	HandleWebSocket(c *gin.Context)
}

type webSocketHandler struct {
	websocketService services.WebSocketService
}

func NewWebSocketHandler(websocketService services.WebSocketService) *webSocketHandler {
	return &webSocketHandler{
		websocketService: websocketService,
	}
}

func (h *webSocketHandler) HandleWebSocket(c *gin.Context) {

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
	h.websocketService.HandleWebSocket(c.Writer, c.Request, saId)
}
