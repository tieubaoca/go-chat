package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tieubaoca/go-chat-server/utils/log"

	"github.com/tieubaoca/go-chat-server/dto/request"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"
)

type AuthenticationHandler interface {
	Login(c *gin.Context)
	logout(c *gin.Context)
}

type authenticationHandler struct {
	websocketService services.WebSocketService
	saasService      services.SaasService
}

func NewAuthenticationHandler(
	websocketService services.WebSocketService,
	saasService services.SaasService,
) *authenticationHandler {
	return &authenticationHandler{
		websocketService: websocketService,
		saasService:      saasService,
	}
}

func (h *authenticationHandler) Login(c *gin.Context) {
	var getAccTokenReq request.GetAccessTokenReq
	err := c.ShouldBindJSON(&getAccTokenReq)
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
	accessToken, refreshToken, err := h.saasService.GetSaasAccessToken(getAccTokenReq.Username, getAccTokenReq.Password)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, response.ResponseData{
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
		Data: map[string]string{
			"access-token":  accessToken,
			"refresh-token": refreshToken,
		},
	})
}

func (h *authenticationHandler) Logout(c *gin.Context) {

	c.SetCookie(
		"access-token",
		"",
		-1,
		"/",
		os.Getenv("DOMAIN"),
		false,
		true,
	)
	c.SetCookie(
		"refresh-token",
		"",
		-1,
		"/",
		os.Getenv("DOMAIN"),
		false,
		true,
	)
	saId, err := h.saasService.GetSaId(utils.GetAccessTokenByReq(c.Request))
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	h.websocketService.Logout(saId)
	c.JSON(
		http.StatusOK,
		response.ResponseData{
			Status:  types.StatusSuccess,
			Message: "Logout successfully",
			Data:    "",
		},
	)
}
