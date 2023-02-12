package services

import (
	"github.com/tieubaoca/go-chat-server/repositories"

	"github.com/tieubaoca/go-chat-server/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageService interface {
	FindMessagesByChatRoomId(chatRoomId string) ([]models.Message, error)
	PaginationMessagesByChatRoomId(chatRoomId string, limit int64, skip int64) ([]models.Message, error)
	InsertMessage(message models.Message) (*mongo.InsertOneResult, error)
}

type messageService struct {
	messageRepository repositories.MessageRepository
}

func NewMessageService(messageRepository repositories.MessageRepository) MessageService {
	return &messageService{messageRepository}
}

func (s *messageService) FindMessagesByChatRoomId(chatRoomId string) ([]models.Message, error) {
	return s.messageRepository.FindMessagesByChatRoomId(chatRoomId)
}

func (s *messageService) InsertMessage(message models.Message) (*mongo.InsertOneResult, error) {
	return s.messageRepository.InsertMessage(message)
}

func (s *messageService) PaginationMessagesByChatRoomId(chatRoomId string, limit int64, skip int64) ([]models.Message, error) {
	return s.messageRepository.PaginationMessagesByChatRoomId(chatRoomId, limit, skip)
}
