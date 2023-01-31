package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/utils"
)

func JwtMiddleware(c *gin.Context) {

	accessTokenString := utils.GetAccessTokenByReq(c.Request)
	if accessTokenString == "" {
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  "Error",
			Message: "Unauthorized",
			Data:    "",
		})
		c.Abort()
		return
	}
	_, err := utils.Parse(accessTokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  "Error",
			Message: "Unauthorized",
			Data:    "",
		})
		c.Abort()
		return
	}
	c.Next()
}
