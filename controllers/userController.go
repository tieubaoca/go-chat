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

// func FindUserByUsername(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	username, ok := vars["username"]
// 	log.Println(username)
// 	if !ok {
// 		w.WriteHeader(http.StatusBadRequest)
// 		response.Res(w, types.StatusError, nil, "Username is empty")
// 		return
// 	}
// 	user, err := services.FindUserByUsername(username)
// 	if err != nil {
// 		w.WriteHeader(http.StatusNoContent)
// 		response.Res(w, types.StatusError, nil, err.Error())
// 		return
// 	}
// 	response.Res(w, types.StatusSuccess, user, "Find user by username successfully")
// }

// func FindUserById(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, ok := vars["id"]
// 	if !ok {
// 		w.WriteHeader(http.StatusBadRequest)
// 		response.Res(w, types.StatusError, nil, "Id is empty")
// 		return
// 	}
// 	user, err := services.FindUserById(id)
// 	if err != nil {
// 		w.WriteHeader(http.StatusNoContent)
// 		response.Res(w, types.StatusError, nil, err.Error())
// 		return
// 	}
// 	response.Res(w, types.StatusSuccess, user, "Find user by id successfully")
// }

func Logout(c *gin.Context) {

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
	services.Logout(utils.GetSaIdFromToken(token), utils.GetSessionIdFromToken(token))
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
	token, err := utils.ParseUnverified(tokenString)
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
	saId := utils.GetSaIdFromToken(token)

	var paginationReq request.PaginationOnlineFriendReq
	err = json.NewDecoder(c.Request.Body).Decode(&paginationReq)
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

	users, err := services.PaginationOnlineFriends(saId, paginationReq)
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
