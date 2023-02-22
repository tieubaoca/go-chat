package services

import (
	"errors"

	"github.com/tieubaoca/go-chat-server/dto/request"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"github.com/tieubaoca/go-chat-server/repositories"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"

	"github.com/tieubaoca/go-chat-server/models"
)

type ChatRoomService interface {
	FindChatRoomById(requester, chatRoomId string) (*response.ChatRoomResponse, error)
	FindChatRoomsBySaId(saId string) ([]response.ChatRoomResponse, error)
	PaginationChatRoomBySaId(req request.ChatRoomPagination) ([]response.ChatRoomResponse, error)
	// CreateNewGroup(requester string, chatRoom request.CreateNewGroupReq) (*mongo.InsertOneResult, error)
	// AddMembersToChatRoom(requester string, req request.AddMemReq) (*mongo.UpdateResult, error)
	// RemoveMembersFromChatRoom(requester string, req request.RemoveMemReq) (*mongo.UpdateResult, error)
	FindDMByMembers(members []string) (*response.ChatRoomResponse, error)
	FindGroupsChatByMembers(members []string) ([]response.ChatRoomResponse, error)
	// CreateNewDMChat(saasToken, member string) (*mongo.InsertOneResult, error)
	// LeaveGroup(requester string, chatRoomId string) (*mongo.UpdateResult, error)
	TransferOwner(requester, chatRoomId, newOwner string) error
}

type chatRoomService struct {
	chatRoomRepository repositories.ChatRoomRepository
	userRepository     repositories.UserRepository
	saasCitizenRepo    repositories.CitizenRepository
}

func NewChatRoomService(
	chatRoomRepository repositories.ChatRoomRepository,
	userRepository repositories.UserRepository,
	saasCitizenRepository repositories.CitizenRepository,
) *chatRoomService {
	return &chatRoomService{
		chatRoomRepository,
		userRepository,
		saasCitizenRepository,
	}
}

func (s *chatRoomService) FindChatRoomById(requester, chatRoomId string) (*response.ChatRoomResponse, error) {
	chatRoom, err := s.chatRoomRepository.FindChatRoomById(chatRoomId)
	if err != nil {
		return nil, err
	}
	if !utils.ContainsString(chatRoom.Members, requester) {
		return nil, errors.New(types.ErrorNotRoomMember)
	}
	return s.convertToChatRoomResponse(chatRoom)
}

func (s *chatRoomService) FindChatRoomsBySaId(saId string) ([]response.ChatRoomResponse, error) {
	chatRooms, err := s.chatRoomRepository.FindChatRoomBySaId(saId)
	if err != nil {
		return nil, err
	}

	return s.convertToChatRoomsResponse(chatRooms)
}

func (s *chatRoomService) FindGroupsChatByMembers(saIds []string) ([]response.ChatRoomResponse, error) {
	chatroms, err := s.chatRoomRepository.FindGroupChatByMembers(saIds)
	if err != nil {
		return nil, err
	}

	return s.convertToChatRoomsResponse(chatroms)
}

func (s *chatRoomService) FindDMByMembers(saIds []string) (*response.ChatRoomResponse, error) {
	chatRoom, err := s.chatRoomRepository.FindDMByMembers(saIds)
	if err != nil {
		return nil, err
	}
	return s.convertToChatRoomResponse(chatRoom)
}
func (s *chatRoomService) convertToChatRoomResponse(chatRoom *models.ChatRoom) (*response.ChatRoomResponse, error) {
	members, err := s.saasCitizenRepo.FindCitizenInList(chatRoom.Members)
	if err != nil {
		return nil, err
	}
	return &response.ChatRoomResponse{
		Id:          chatRoom.Id,
		Name:        chatRoom.Name,
		Owner:       chatRoom.Owner,
		Type:        chatRoom.Type,
		LastMessage: chatRoom.LastMessage,
		IsBlocked:   chatRoom.IsBlocked,
		Members:     members,
	}, nil
}

func (s *chatRoomService) convertToChatRoomsResponse(chatRooms []models.ChatRoom) ([]response.ChatRoomResponse, error) {
	var chatRoomRes []response.ChatRoomResponse
	for _, room := range chatRooms {
		chatRoom, err := s.convertToChatRoomResponse(&room)
		if err != nil {
			return nil, err
		}
		chatRoomRes = append(chatRoomRes, *chatRoom)
	}
	return chatRoomRes, nil
}

// func (s *chatRoomService) CreateNewGroup(requester string, req request.CreateNewGroupReq) (*mongo.InsertOneResult, error) {
// 	if !utils.ContainsString(req.Members, requester) {
// 		req.Members = append(req.Members, requester)
// 	}
// 	for _, member := range req.Members {
// 		if !s.userRepository.IsUserExist(member) {
// 			return nil, errors.New(types.ErrorUserNotExist)
// 		}
// 	}
// 	chatRoom := models.ChatRoom{
// 		Owner:     requester,
// 		Type:      models.ChatRoomTypeGroup,
// 		Name:      req.Name,
// 		Members:   req.Members,
// 		IsBlocked: false,
// 	}
// 	return s.chatRoomRepository.InsertChatRoom(chatRoom)
// }

// func (s *chatRoomService) CreateNewDMChat(saasToken, member string) (*mongo.InsertOneResult, error) {
// 	friends, err := utils.GetAllFriends(saasToken)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, friend := range friends {
// 		if friend.(map[string]interface{})["saIdFriend"].(string) == member {
// 			saId, err := utils.GetSaIdFromToken(saasToken)
// 			if err != nil {
// 				return nil, err
// 			}
// 			members := []string{
// 				member,
// 				saId,
// 			}
// 			sort.Strings(members)
// 			return s.chatRoomRepository.InsertChatRoom(
// 				models.ChatRoom{
// 					Owner:     "",
// 					Type:      models.ChatRoomTypeDM,
// 					Name:      members[0] + "-" + members[1],
// 					Members:   members,
// 					IsBlocked: false,
// 				},
// 			)
// 		}
// 	}
// 	return nil, errors.New("only dm friend")

// }

// func (s *chatRoomService) AddMembersToChatRoom(requester string, req request.AddMemReq) (*mongo.UpdateResult, error) {
// 	chatRoom, err := s.FindChatRoomById(requester, req.ChatRoomId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if chatRoom.Type == models.ChatRoomTypeDM {
// 		return nil, errors.New(types.ErrorInvalidChatRoomType)
// 	}
// 	if chatRoom.Owner != requester {
// 		return nil, errors.New(types.ErrorOnlyOwner)
// 	}

// 	return s.chatRoomRepository.AddMembersToChatRoom(req.ChatRoomId, req.SaIds)
// }

// func (s *chatRoomService) RemoveMembersFromChatRoom(requester string, req request.RemoveMemReq) (*mongo.UpdateResult, error) {
// 	chatRoom, err := s.FindChatRoomById(requester, req.ChatRoomId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if chatRoom.Type == models.ChatRoomTypeDM {
// 		return nil, errors.New(types.ErrorInvalidChatRoomType)
// 	}
// 	if chatRoom.Owner != requester {
// 		return nil, errors.New(types.ErrorOnlyOwner)
// 	}
// 	if !utils.ContainsString(chatRoom.Members, requester) {
// 		return nil, errors.New(types.ErrorInvalidInput)
// 	}
// 	return s.chatRoomRepository.RemoveMembersFromChatRoom(req.ChatRoomId, req.SaIds)
// }

// func (s *chatRoomService) LeaveGroup(requester string, chatRoomId string) (*mongo.UpdateResult, error) {
// 	chatRoom, err := s.FindChatRoomById(requester, chatRoomId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if chatRoom.Type == models.ChatRoomTypeDM {
// 		return nil, errors.New(types.ErrorInvalidChatRoomType)
// 	}
// 	if !utils.ContainsString(chatRoom.Members, requester) {
// 		return nil, errors.New(types.ErrorInvalidInput)
// 	}
// 	if chatRoom.Owner == requester {
// 		return nil, errors.New(types.ErrorInvalidInput)
// 	}

// 	return s.chatRoomRepository.RemoveMembersFromChatRoom(chatRoomId, []string{requester})
// }

func (s *chatRoomService) TransferOwner(requester, chatRoomId, newOwner string) error {
	chatRoom, err := s.FindChatRoomById(requester, chatRoomId)
	if err != nil {
		return err
	}
	if chatRoom.Owner != requester {
		return errors.New(types.ErrorOnlyOwner)
	}
	if !s.userRepository.IsUserExist(newOwner) {
		return errors.New(types.ErrorUserNotExist)
	}
	return s.chatRoomRepository.TransferOwner(chatRoomId, newOwner)
}

func (s *chatRoomService) PaginationChatRoomBySaId(req request.ChatRoomPagination) ([]response.ChatRoomResponse, error) {
	chatRooms, err := s.chatRoomRepository.PaginationChatRoomBySaId(
		req.SaId,
		req.Skip,
		req.Limit,
	)
	if err != nil {
		return nil, err
	}
	return s.convertToChatRoomsResponse(chatRooms)
}
