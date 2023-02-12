package services

import (
	"github.com/tieubaoca/go-chat-server/repositories"

	"github.com/tieubaoca/go-chat-server/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRoomService interface {
	FindChatRoomById(id string) (models.ChatRoom, error)
	FindChatRoomsBySaId(saId string) ([]models.ChatRoom, error)
	InsertChatRoom(chatRoom models.ChatRoom) (*mongo.InsertOneResult, error)
	AddMembersToChatRoom(chatRoomId string, members []string) (*mongo.UpdateResult, error)
	RemoveMembersFromChatRoom(chatRoomId string, members []string) (*mongo.UpdateResult, error)
	FindDMByMembers(members []string) (models.ChatRoom, error)
	FindGroupsChatByMembers(members []string) ([]models.ChatRoom, error)
}

type chatRoomService struct {
	chatRoomRepository repositories.ChatRoomRepository
}

func NewChatRoomService(chatRoomRepository repositories.ChatRoomRepository) ChatRoomService {
	return &chatRoomService{chatRoomRepository}
}

func (s *chatRoomService) FindChatRoomById(id string) (models.ChatRoom, error) {
	return s.chatRoomRepository.FindChatRoomById(id)
}

func (s *chatRoomService) FindChatRoomsBySaId(saId string) ([]models.ChatRoom, error) {
	return s.chatRoomRepository.FindChatRoomBySaId(saId)
}

func (s *chatRoomService) FindGroupsChatByMembers(saIds []string) ([]models.ChatRoom, error) {
	return s.chatRoomRepository.FindGroupChatByMembers(saIds)
}

func (s *chatRoomService) FindDMByMembers(saIds []string) (models.ChatRoom, error) {
	return s.chatRoomRepository.FindDMByMembers(saIds)
}

func (s *chatRoomService) InsertChatRoom(chatRoom models.ChatRoom) (*mongo.InsertOneResult, error) {
	return s.chatRoomRepository.InsertChatRoom(chatRoom)
}

func (s *chatRoomService) AddMembersToChatRoom(chatRoomId string, members []string) (*mongo.UpdateResult, error) {
	return s.chatRoomRepository.AddMembersToChatRoom(chatRoomId, members)
}

func (s *chatRoomService) RemoveMembersFromChatRoom(chatRoomId string, members []string) (*mongo.UpdateResult, error) {
	return s.chatRoomRepository.RemoveMembersFromChatRoom(chatRoomId, members)
}
