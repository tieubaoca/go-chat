package controllers

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/models"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"
	"github.com/tieubaoca/go-chat-server/utils/log"
)

func FindChatRoomById(c *gin.Context) {
	vars := mux.Vars(c.Request)
	id, ok := vars["id"]
	if !ok {
		log.ErrorLogger.Println(types.ErrorInvalidInput)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorInvalidInput,
			Data:    "",
		})
		return
	}
	chatRoom, err := services.FindChatRoomById(id)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusNoContent, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	token, _ := utils.ParseUnverified(utils.GetAccessTokenByReq(c.Request))
	if !utils.ContainsString(chatRoom.Members, utils.GetSaIdFromToken(token)) {
		log.ErrorLogger.Println(types.ErrorNotRoomMember)
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorNotRoomMember,
			Data:    "",
		})
		return
	}
	c.JSON(http.StatusOK, response.ResponseData{
		Status:  types.StatusSuccess,
		Message: "OK",
		Data:    chatRoom,
	})
}

func FindChatRooms(c *gin.Context) {
	token, err := utils.ParseUnverified(utils.GetAccessTokenByReq(c.Request))
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	chatRooms, err := services.FindChatRoomsByMember(utils.GetSaIdFromToken(token))
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
		Message: "OK",
		Data:    chatRooms,
	})
}

func FindDMByMembers(c *gin.Context) {
	var member string
	err := json.NewDecoder(c.Request.Body).Decode(&member)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorInvalidInput,
			Data:    "",
		})
		return
	}
	if member == "" {
		log.ErrorLogger.Println(types.ErrorInvalidInput)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorInvalidInput,
			Data:    "",
		})
		return
	}
	token, err := utils.ParseUnverified(utils.GetAccessTokenByReq(c.Request))
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}

	members := []string{member, utils.GetSaIdFromToken(token)}

	chatRoom, err := services.FindDMByMembers(members)
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
		Message: "OK",
		Data:    chatRoom,
	})

}

func FindGroupsByMembers(c *gin.Context) {
	var members []string
	err := json.NewDecoder(c.Request.Body).Decode(&members)
	if err != nil {
		members = make([]string, 0)
	}
	token, err := utils.ParseUnverified(utils.GetAccessTokenByReq(c.Request))
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}

	chatRooms, err := services.FindGroupsByMembers(append(members, utils.GetSaIdFromToken(token)))
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
		Message: "OK",
		Data:    chatRooms,
	})
}

func CreateNewGroupChat(c *gin.Context) {

	var chatRoom models.ChatRoom
	err := json.NewDecoder(c.Request.Body).Decode(&chatRoom)

	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorInvalidInput,
			Data:    "",
		})
		return
	}
	chatRoom.Type = models.ChatRoomTypeGroup
	token, _ := utils.ParseUnverified(utils.GetAccessTokenByReq(c.Request))
	chatRoom.Owner = utils.GetSaIdFromToken(token)

	result, err := services.InsertChatRoom(chatRoom)
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
		Message: "OK",
		Data:    result,
	})
}

func CreateDMRoom(c *gin.Context) {
	var member string
	err := json.NewDecoder(c.Request.Body).Decode(&member)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorInvalidInput,
			Data:    "",
		})
		return
	}
	if member == "" {
		log.ErrorLogger.Println(types.ErrorInvalidInput)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorInvalidInput,
			Data:    "",
		})
		return
	}
	token, err := utils.ParseUnverified(utils.GetAccessTokenByReq(c.Request))
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}

	members := []string{
		member,
		utils.GetSaIdFromToken(token),
	}
	sort.Strings(members)
	result, err := services.InsertChatRoom(
		models.ChatRoom{
			Name:    members[0] + "-" + members[1],
			Type:    models.ChatRoomTypeDM,
			Members: members,
		},
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
		Message: "OK",
		Data:    result,
	})
}
