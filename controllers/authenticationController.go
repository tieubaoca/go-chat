package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tieubaoca/go-chat-server/utils/log"

	"github.com/tieubaoca/go-chat-server/dto/request"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"
)

func Authentication(c *gin.Context) {
	tokenString := utils.GetAccessTokenByReq(c.Request)
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorTokenEmpty,
			Data:    "",
		})
		return
	}
	token, err := utils.JWTParseUnverified(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	c.JSON(http.StatusOK, response.ResponseData{
		Status:  types.StatusSuccess,
		Message: "OK",
		Data:    token.Claims,
	})
}

func GetAccessToken(c *gin.Context) {
	var getAccTokenReq request.GetAccessTokenReq
	err := json.NewDecoder(c.Request.Body).Decode(&getAccTokenReq)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorInvalidInput,
			Data:    "",
		})
		return
	}
	if getAccTokenReq.Username == "" || getAccTokenReq.Password == "" {
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorInvalidInput,
			Data:    "",
		})
		return
	}
	accessToken, refreshToken, err := utils.GetSaasAccessToken(getAccTokenReq.Username, getAccTokenReq.Password)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}

	c.SetCookie(
		"access-token",
		accessToken,
		10*60,
		"/",
		os.Getenv("DOMAIN"),
		false,
		false,
	)
	c.SetCookie(
		"refresh-token",
		refreshToken,
		10*60,
		"/",
		os.Getenv("DOMAIN"),
		false,
		false,
	)
	_, err = utils.JWTSaasParse(accessToken)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}

	c.JSON(http.StatusOK, response.ResponseData{
		Status:  types.StatusSuccess,
		Message: "OK",
		Data: map[string]string{
			"access-token":  accessToken,
			"refresh-token": refreshToken,
		},
	})

}
