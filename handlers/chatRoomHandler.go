package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tieubaoca/go-chat-server/dto/request"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"
	"github.com/tieubaoca/go-chat-server/utils/log"
)

type ChatRoomHandler interface {
	FindChatRoomById(c *gin.Context)
	FindChatRoomsBySaId(c *gin.Context)
	FindDMByMember(c *gin.Context)
	PaginationChatRoomBySaId(c *gin.Context)
	// CreateNewGroupChat(c *gin.Context)
	// CreateNewDMChat(c *gin.Context)
	// AddMembersToGroup(c *gin.Context)
	// RemoveMembersFromGroup(c *gin.Context)
	// LeaveGroup(c *gin.Context)
}

type chatRoomHandler struct {
	chatRoomService services.ChatRoomService
	saasService     services.SaasService
}

func NewChatRoomHandler(
	chatRoomService services.ChatRoomService,
	saasService services.SaasService,
) *chatRoomHandler {
	return &chatRoomHandler{
		chatRoomService: chatRoomService,
		saasService:     saasService,
	}
}

func (h *chatRoomHandler) FindChatRoomById(c *gin.Context) {
	id := c.Param("id")
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
	log.InfoLogger.Println("saId: ", saId)
	chatRoom, err := h.chatRoomService.FindChatRoomById(
		saId,
		id,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			response.ResponseData{
				Status:  types.StatusError,
				Message: err.Error(),
				Data:    nil,
			},
		)
		return
	}
	c.JSON(http.StatusOK, response.ResponseData{
		Status:  types.StatusSuccess,
		Message: "OK",
		Data:    chatRoom,
	})
}

func (h *chatRoomHandler) FindChatRoomsBySaId(c *gin.Context) {
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

func (h *chatRoomHandler) PaginationChatRoomBySaId(c *gin.Context) {
	var req request.ChatRoomPagination
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
	saId, err := h.saasService.GetSaId(utils.GetAccessTokenByReq(c.Request))
	req.SaId = saId
	if err != nil {
		log.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, response.ResponseData{
			Status:  types.StatusError,
			Message: err.Error(),
			Data:    "",
		})
		return
	}
	chatRooms, err := h.chatRoomService.PaginationChatRoomBySaId(
		req,
	)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			response.ResponseData{
				Status:  types.StatusError,
				Message: err.Error(),
				Data:    nil,
			},
		)
		return
	}
	c.JSON(http.StatusOK, response.ResponseData{
		Status:  types.StatusSuccess,
		Message: "OK",
		Data:    chatRooms,
	})
}

// func (h *chatRoomHandler) CreateNewGroupChat(c *gin.Context) {

// 	var req request.CreateNewGroupReq
// 	err := c.ShouldBindJSON(&req)
// 	if err != nil {
// 		log.ErrorLogger.Println(err)
// 		c.JSON(http.StatusBadRequest, response.ResponseData{
// 			Status:  types.StatusError,
// 			Message: types.ErrorInvalidInput,
// 			Data:    "",
// 		})
// 		return
// 	}
// 	saId, err := h.saasService.GetSaId(utils.GetAccessTokenByReq(c.Request))
// 	if err != nil {
// 		log.ErrorLogger.Println(err)
// 		c.JSON(http.StatusInternalServerError, response.ResponseData{
// 			Status:  types.StatusError,
// 			Message: err.Error(),
// 			Data:    "",
// 		})
// 	}

// 	result, err := h.chatRoomService.CreateNewGroup(saId, req)
// 	if err != nil {
// 		log.ErrorLogger.Println(err)
// 		c.JSON(http.StatusInternalServerError, response.ResponseData{
// 			Status:  types.StatusError,
// 			Message: err.Error(),
// 			Data:    "",
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, response.ResponseData{
// 		Status:  types.StatusSuccess,
// 		Message: "OK",
// 		Data:    result,
// 	})
// }

// func (h *chatRoomHandler) CreateNewDMChat(c *gin.Context) {
// 	var member string
// 	err := c.ShouldBindJSON(&member)
// 	if err != nil {
// 		log.ErrorLogger.Println(err)
// 		c.JSON(http.StatusBadRequest, response.ResponseData{
// 			Status:  types.StatusError,
// 			Message: types.ErrorInvalidInput,
// 			Data:    "",
// 		})
// 		return
// 	}
// 	if member == "" {
// 		log.ErrorLogger.Println(types.ErrorInvalidInput)
// 		c.JSON(http.StatusBadRequest, response.ResponseData{
// 			Status:  types.StatusError,
// 			Message: types.ErrorInvalidInput,
// 			Data:    "",
// 		})
// 		return
// 	}

// 	result, err := h.chatRoomService.CreateNewDMChat(
// 		utils.GetAccessTokenByReq(c.Request),
// 		member,
// 	)
// 	if err != nil {
// 		log.ErrorLogger.Println(err)
// 		c.JSON(http.StatusInternalServerError, response.ResponseData{
// 			Status:  types.StatusError,
// 			Message: err.Error(),
// 			Data:    "",
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, response.ResponseData{
// 		Status:  types.StatusSuccess,
// 		Message: "OK",
// 		Data:    result,
// 	})
// }

// func (h *chatRoomHandler) AddMemberToGroup(c *gin.Context) {
// 	var req request.AddMemReq
// 	err := c.ShouldBindJSON(&req)
// 	if err != nil {
// 		log.ErrorLogger.Println(err)
// 		c.JSON(http.StatusBadRequest, response.ResponseData{
// 			Status:  types.StatusError,
// 			Message: types.ErrorInvalidInput,
// 			Data:    "",
// 		})
// 		return
// 	}
// 	saId, err := h.saasService.GetSaId(utils.GetAccessTokenByReq(c.Request))
// 	if err != nil {
// 		log.ErrorLogger.Println(err)
// 		c.JSON(http.StatusInternalServerError, response.ResponseData{
// 			Status:  types.StatusError,
// 			Message: err.Error(),
// 			Data:    "",
// 		})
// 		return
// 	}
// 	result, err := h.chatRoomService.AddMembersToChatRoom(saId, req)
// 	if err != nil {
// 		log.ErrorLogger.Println(err)
// 		c.JSON(http.StatusInternalServerError, response.ResponseData{
// 			Status:  types.StatusError,
// 			Message: err.Error(),
// 			Data:    "",
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, response.ResponseData{
// 		Status:  types.StatusSuccess,
// 		Message: "OK",
// 		Data:    result,
// 	})
// }

// func (h *chatRoomHandler) RemoveMemberFromGroup(c *gin.Context) {
// 	var req request.RemoveMemReq
// 	err := c.BindJSON(&req)
// 	if err != nil {
// 		log.ErrorLogger.Println(err)
// 		c.JSON(http.StatusBadRequest, response.ResponseData{
// 			Status:  types.StatusError,
// 			Message: types.ErrorInvalidInput,
// 			Data:    "",
// 		})
// 		return
// 	}
// 	saId, err := h.saasService.GetSaId(utils.GetAccessTokenByReq(c.Request))
// 	if err != nil {
// 		log.ErrorLogger.Println(err)
// 		c.JSON(http.StatusInternalServerError, response.ResponseData{
// 			Status:  types.StatusError,
// 			Message: err.Error(),
// 			Data:    "",
// 		})
// 		return
// 	}
// 	result, err := h.chatRoomService.RemoveMembersFromChatRoom(saId, req)
// 	if err != nil {
// 		log.ErrorLogger.Println(err)
// 		c.JSON(http.StatusInternalServerError, response.ResponseData{
// 			Status:  types.StatusError,
// 			Message: err.Error(),
// 			Data:    "",
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, response.ResponseData{
// 		Status:  types.StatusSuccess,
// 		Message: "OK",
// 		Data:    result,
// 	})
// }

// func (h *chatRoomHandler) LeaveGroup(c *gin.Context) {
// 	chatRoomId := c.Param("chatRoomId")
// 	saId, err := h.saasService.GetSaId(utils.GetAccessTokenByReq(c.Request))
// 	if err != nil {
// 		log.ErrorLogger.Println(err)
// 		c.JSON(http.StatusInternalServerError, response.ResponseData{
// 			Status:  types.StatusError,
// 			Message: err.Error(),
// 			Data:    "",
// 		})
// 		return
// 	}
// 	result, err := h.chatRoomService.LeaveGroup(saId, chatRoomId)
// 	if err != nil {
// 		log.ErrorLogger.Println(err)
// 		c.JSON(http.StatusInternalServerError, response.ResponseData{
// 			Status:  types.StatusError,
// 			Message: err.Error(),
// 			Data:    "",
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, response.ResponseData{
// 		Status:  types.StatusSuccess,
// 		Message: "OK",
// 		Data:    result,
// 	})
// }
