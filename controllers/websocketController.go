package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils/log"

	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/utils"
)

func HandleWebSocket(c *gin.Context) {

	token, err := utils.ParseUnverified(utils.GetAccessTokenByReq(c.Request))
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(
			http.StatusUnauthorized,
			response.ResponseData{
				Status:  types.StatusError,
				Message: err.Error(),
				Data:    "",
			},
		)
		return
	}

	services.HandleWebSocket(c.Writer, c.Request, utils.GetSaIdFromToken(token))
}
