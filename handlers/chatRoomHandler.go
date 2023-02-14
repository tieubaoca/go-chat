package handlers

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

type ChatRoomHandler interface {
	FindChatRoomById(c *gin.Context)
	FindChatRoomsBySaId(c *gin.Context)
	FindDMByMember(c *gin.Context)
	CreateNewGroupChat(c *gin.Context)
	CreateNewDMChat(c *gin.Context)
	AddMembersToGroup(c *gin.Context)
	RemoveMembersFromGroup(c *gin.Context)
	LeaveGroup(c *gin.Context)
}

type chatRoomHandler struct {
	chatRoomService services.ChatRoomService
}

func NewChatRoomHandler(chatRoomService services.ChatRoomService) *chatRoomHandler {
	return &chatRoomHandler{chatRoomService: chatRoomService}
}

func (h *chatRoomHandler) FindChatRoomById(c *gin.Context) {
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
	chatRoom, err := h.chatRoomService.FindChatRoomById(id)
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

func (h *chatRoomHandler) FindChatRoomsBySaId(c *gin.Context) {
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
	chatRooms, err := h.chatRoomService.FindChatRoomsBySaId(saId)
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

func (h *chatRoomHandler) FindDMByMember(c *gin.Context) {
	member := c.Param("member")
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

	chatRoom, err := h.chatRoomService.FindDMByMembers(members)
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

func (h *chatRoomHandler) FindGroupsByMembers(c *gin.Context) {
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

	chatRooms, err := h.chatRoomService.FindGroupsChatByMembers(append(members, saId))
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

func (h *chatRoomHandler) CreateNewGroupChat(c *gin.Context) {

	var chatRoom models.ChatRoom
	err := c.ShouldBindJSON(&chatRoom)
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
	chatRoom.Owner = saId
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
		chatRoom.Members = append(chatRoom.Members, saId)
	}

	result, err := h.chatRoomService.InsertChatRoom(chatRoom)
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

func (h *chatRoomHandler) CreateNewDMChat(c *gin.Context) {
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
		friendIds = append(friendIds, friend.(map[string]interface{})["saIdFriend"].(string))
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
	result, err := h.chatRoomService.InsertChatRoom(
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

func (h *chatRoomHandler) AddMemberToGroup(c *gin.Context) {
	var req request.AddMemReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorInvalidInput,
			Data:    "",
		})
		return
	}
	chatRoom, err := h.chatRoomService.FindChatRoomById(req.ChatRoomId)
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
	result, err := h.chatRoomService.AddMembersToChatRoom(req.ChatRoomId, req.SaIds)
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

func (h *chatRoomHandler) RemoveMemberFromGroup(c *gin.Context) {
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
	chatRoom, err := h.chatRoomService.FindChatRoomById(req.ChatRoomId)
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
	if len(req.SaIds) == 0 {
		log.ErrorLogger.Println(types.ErrorInvalidInput)
		c.JSON(http.StatusBadRequest, response.ResponseData{
			Status:  types.StatusError,
			Message: types.ErrorInvalidInput,
			Data:    "",
		})
		return
	}
	if utils.ContainsString(req.SaIds, saId) {
		c.JSON(http.StatusUnauthorized, response.ResponseData{
			Status:  types.StatusError,
			Message: "You can't remove yourself from the group",
			Data:    "",
		})
		return
	}
	result, err := h.chatRoomService.RemoveMembersFromChatRoom(req.ChatRoomId, req.SaIds)
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

func (h *chatRoomHandler) LeaveGroup(c *gin.Context) {
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
	chatRoom, err := h.chatRoomService.FindChatRoomById(chatRoomId)
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
	if chatRoom.Type == models.ChatRoomTypeDM {
		c.JSON(http.StatusInternalServerError, response.ResponseData{
			Status:  types.StatusError,
			Message: "You can't leave DM",
			Data:    "",
		},
		)
		return
	}
	if !utils.ContainsString(chatRoom.Members, saId) {
		log.ErrorLogger.Println(types.ErrorUnauthorized)
		c.JSON(http.StatusInternalServerError, response.ResponseData{
			Status:  types.StatusError,
			Message: "You are not member of this group",
			Data:    "",
		})
		return
	}
	if chatRoom.Owner == saId {
		c.JSON(http.StatusInternalServerError, response.ResponseData{
			Status:  types.StatusError,
			Message: "You are owner of this group",
			Data:    "",
		})
		return
	}
	result, err := h.chatRoomService.RemoveMembersFromChatRoom(chatRoomId, []string{saId})
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
