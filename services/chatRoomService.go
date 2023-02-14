package services

import (
	"errors"
	"sort"

	"github.com/tieubaoca/go-chat-server/dto/request"
	"github.com/tieubaoca/go-chat-server/repositories"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"

	"github.com/tieubaoca/go-chat-server/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRoomService interface {
	FindChatRoomById(requester, chatRoomId string) (*models.ChatRoom, error)
	FindChatRoomsBySaId(saId string) ([]models.ChatRoom, error)
	CreateNewGroup(requester string, chatRoom request.CreateNewGroupReq) (*mongo.InsertOneResult, error)
	AddMembersToChatRoom(requester string, req request.AddMemReq) (*mongo.UpdateResult, error)
	RemoveMembersFromChatRoom(requester string, req request.RemoveMemReq) (*mongo.UpdateResult, error)
	FindDMByMembers(members []string) (*models.ChatRoom, error)
	FindGroupsChatByMembers(members []string) ([]models.ChatRoom, error)
	CreateNewDMChat(saasToken, member string) (*mongo.InsertOneResult, error)
	LeaveGroup(requester string, chatRoomId string) (*mongo.UpdateResult, error)
}

type chatRoomService struct {
	chatRoomRepository repositories.ChatRoomRepository
}

func NewChatRoomService(chatRoomRepository repositories.ChatRoomRepository) *chatRoomService {
	return &chatRoomService{chatRoomRepository}
}

func (s *chatRoomService) FindChatRoomById(requester, chatRoomId string) (*models.ChatRoom, error) {
	chatRoom, err := s.chatRoomRepository.FindChatRoomById(chatRoomId)
	if err != nil {
		return nil, err
	}
	if !utils.ContainsString(chatRoom.Members, requester) {
		return nil, errors.New(types.ErrorNotRoomMember)
	}
	return chatRoom, nil
}

func (s *chatRoomService) FindChatRoomsBySaId(saId string) ([]models.ChatRoom, error) {
	return s.chatRoomRepository.FindChatRoomBySaId(saId)
}

func (s *chatRoomService) FindGroupsChatByMembers(saIds []string) ([]models.ChatRoom, error) {
	return s.chatRoomRepository.FindGroupChatByMembers(saIds)
}

func (s *chatRoomService) FindDMByMembers(saIds []string) (*models.ChatRoom, error) {
	return s.chatRoomRepository.FindDMByMembers(saIds)
}

func (s *chatRoomService) CreateNewGroup(requester string, req request.CreateNewGroupReq) (*mongo.InsertOneResult, error) {
	if !utils.ContainsString(req.Members, requester) {
		req.Members = append(req.Members, requester)
	}
	chatRoom := models.ChatRoom{
		Owner:     requester,
		Type:      models.ChatRoomTypeGroup,
		Name:      req.Name,
		Members:   req.Members,
		IsBlocked: false,
	}
	return s.chatRoomRepository.InsertChatRoom(chatRoom)
}

func (s *chatRoomService) CreateNewDMChat(saasToken, member string) (*mongo.InsertOneResult, error) {
	friends, err := utils.GetAllFriends(saasToken)
	if err != nil {
		return nil, err
	}
	for _, friend := range friends {
		if friend.(map[string]interface{})["saIdFriend"].(string) == member {
			saId, err := utils.GetSaIdFromToken(saasToken)
			if err != nil {
				return nil, err
			}
			members := []string{
				member,
				saId,
			}
			sort.Strings(members)
			return s.chatRoomRepository.InsertChatRoom(
				models.ChatRoom{
					Owner:     "",
					Type:      models.ChatRoomTypeDM,
					Name:      members[0] + "-" + members[1],
					Members:   members,
					IsBlocked: false,
				},
			)
		}
	}
	return nil, errors.New("only dm friend")

}

func (s *chatRoomService) AddMembersToChatRoom(requester string, req request.AddMemReq) (*mongo.UpdateResult, error) {
	chatRoom, err := s.FindChatRoomById(requester, req.ChatRoomId)
	if err != nil {
		return nil, err
	}
	if chatRoom.Type == models.ChatRoomTypeDM {
		return nil, errors.New(types.ErrorInvalidChatRoomType)
	}
	if chatRoom.Owner != requester {
		return nil, errors.New(types.ErrorOnlyOwner)
	}

	return s.chatRoomRepository.AddMembersToChatRoom(req.ChatRoomId, req.SaIds)
}

func (s *chatRoomService) RemoveMembersFromChatRoom(requester string, req request.RemoveMemReq) (*mongo.UpdateResult, error) {
	chatRoom, err := s.FindChatRoomById(requester, req.ChatRoomId)
	if err != nil {
		return nil, err
	}
	if chatRoom.Type == models.ChatRoomTypeDM {
		return nil, errors.New(types.ErrorInvalidChatRoomType)
	}
	if chatRoom.Owner != requester {
		return nil, errors.New(types.ErrorOnlyOwner)
	}
	if !utils.ContainsString(chatRoom.Members, requester) {
		return nil, errors.New(types.ErrorInvalidInput)
	}
	return s.chatRoomRepository.RemoveMembersFromChatRoom(req.ChatRoomId, req.SaIds)
}

func (s *chatRoomService) LeaveGroup(requester string, chatRoomId string) (*mongo.UpdateResult, error) {
	chatRoom, err := s.FindChatRoomById(requester, chatRoomId)
	if err != nil {
		return nil, err
	}
	if chatRoom.Type == models.ChatRoomTypeDM {
		return nil, errors.New(types.ErrorInvalidChatRoomType)
	}
	if !utils.ContainsString(chatRoom.Members, requester) {
		return nil, errors.New(types.ErrorInvalidInput)
	}
	if chatRoom.Owner == requester {
		return nil, errors.New(types.ErrorInvalidInput)
	}

	return s.chatRoomRepository.RemoveMembersFromChatRoom(chatRoomId, []string{requester})
}
