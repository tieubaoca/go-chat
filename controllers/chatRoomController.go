package controllers

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/tieubaoca/go-chat-server/dto/request"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/models"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"
	"github.com/tieubaoca/go-chat-server/utils/log"
)

func FindChatRoomById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		log.ErrorLogger.Println(types.ErrorInvalidInput)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorInvalidInput,
			Data:    "",
		})
		return
	}
	log.InfoLogger.Println("ID:", id)
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
	if !utils.ContainsString(chatRoom.Members, saId) {
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
	chatRooms, err := services.FindChatRoomsByMember(saId)
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
	// err := json.NewDecoder(c.Request.Body).Decode(&member)
	err := c.ShouldBindJSON(&member)
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

	members := []string{member, saId}

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

	chatRooms, err := services.FindGroupsByMembers(append(members, saId))
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
	saId, err := utils.GetSaIdFromToken(utils.GetAccessTokenByReq(c.Request))
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
	}
	chatRoom.Owner, err = utils.GetSaIdFromToken(saId)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}

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
	friends, err := utils.GetAllFriends(utils.GetAccessTokenByReq(c.Request))
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}

	var friendIds []string
	for _, friend := range friends {
		friendIds = append(friendIds, friend.(map[string]interface{})["id"].(string))
	}

	if !utils.ContainsString(friendIds, member) {
		log.ErrorLogger.Println(types.ErrorInvalidInput)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorInvalidInput,
			Data:    "",
		})
		return
	}

	saId, err := utils.GetSaIdFromToken(utils.GetAccessTokenByReq(c.Request))
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
		saId,
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

func AddMemberToGroup(c *gin.Context) {
	var req request.AddMemReq
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorInvalidInput,
			Data:    "",
		})
		return
	}
	chatRoom, err := services.FindChatRoomById(req.ChatRoomId)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	saId, err := utils.GetSaIdFromToken(utils.GetAccessTokenByReq(c.Request))
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	if chatRoom.Owner != saId {
		log.ErrorLogger.Println(types.ErrorUnauthorized)
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorUnauthorized,
			Data:    "",
		})
		return
	}
	result, err := services.AddMemberToChatRoom(req.ChatRoomId, req.SaIds)
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

func RemoveMemberFromGroup(c *gin.Context) {
	var req request.RemoveMemReq
	err := c.BindJSON(&req)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorInvalidInput,
			Data:    "",
		})
		return
	}
	chatRoom, err := services.FindChatRoomById(req.ChatRoomId)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	saId, err := utils.GetSaIdFromToken(utils.GetAccessTokenByReq(c.Request))
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	if chatRoom.Owner != saId {
		log.ErrorLogger.Println(types.ErrorUnauthorized)
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorUnauthorized,
			Data:    "",
		})
		return
	}
	result, err := services.RemoveMemberFromChatRoom(req.ChatRoomId, req.SaIds)
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

func LeaveGroup(c *gin.Context) {
	chatRoomId := c.Param("chatRoomId")
	if chatRoomId == "" {
		log.ErrorLogger.Println(types.ErrorInvalidInput)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorInvalidInput,
			Data:    "",
		})
		return
	}
	chatRoom, err := services.FindChatRoomById(chatRoomId)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	saId, err := utils.GetSaIdFromToken(utils.GetAccessTokenByReq(c.Request))
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	if utils.ContainsString(chatRoom.Members, saId) {
		log.ErrorLogger.Println(types.ErrorUnauthorized)
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorUnauthorized,
			Data:    "",
		})
		return
	}
	result, err := services.RemoveMemberFromChatRoom(chatRoomId, []string{saId})
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
