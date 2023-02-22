package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/utils/log"

	"github.com/tieubaoca/go-chat-server/dto/request"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"
)

type UserHandler interface {
	PaginationOnlineFriend(c *gin.Context)
}

type userHandler struct {
	userService services.UserService
	saasService services.SaasService
}

func NewUserHandler(
	userService services.UserService,
	saasService services.SaasService,
) *userHandler {
	return &userHandler{
		userService: userService,
		saasService: saasService,
	}
}

func (h *userHandler) PaginationOnlineFriend(c *gin.Context) {
	tokenString := utils.GetAccessTokenByReq(c.Request)
	var paginationReq request.PaginationOnlineFriendReq
	err := c.ShouldBindJSON(&paginationReq)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(
			http.StatusBadRequest,
			response.ResponseData{
				Status:  types.StatusError,
				Message: err.Error(),
				Data:    "",
			},
		)
		return
	}

	users, err := h.saasService.GetListFriendInfo(tokenString, paginationReq)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(
			http.StatusInternalServerError,
			response.ResponseData{
				Status:  types.StatusError,
				Message: err.Error(),
				Data:    "",
			},
		)
		return
	}
	friendStatus, err := h.userService.FindUserStatusInUserList(
		users,
	)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(
			http.StatusInternalServerError,
			response.ResponseData{
				Status:  types.StatusError,
				Message: err.Error(),
				Data:    "",
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		response.ResponseData{
			Status:  types.StatusSuccess,
			Message: "Pagination online friends successfully",
			Data:    friendStatus,
		},
	)
}
