package controllers

import (
	"encoding/json"
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

func Logout(c *gin.Context) {

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
	services.Logout(saId)
	c.JSON(
		http.StatusOK,
		response.ResponseData{
			Status:  types.StatusSuccess,
			Message: "Logout successfully",
			Data:    "",
		},
	)
}

func PaginationOnlineFriend(c *gin.Context) {
	tokenString := utils.GetAccessTokenByReq(c.Request)
	var paginationReq request.PaginationOnlineFriendReq
	err := json.NewDecoder(c.Request.Body).Decode(&paginationReq)
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

	users, err := utils.GetListFriendInfo(tokenString, paginationReq)
	log.InfoLogger.Println(users)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(
			http.StatusNoContent,
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
			Data:    users,
		},
	)
}
