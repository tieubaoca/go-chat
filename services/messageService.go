package services

import (
	"github.com/tieubaoca/go-chat-server/dto/request"
	"github.com/tieubaoca/go-chat-server/repositories"
	"github.com/tieubaoca/go-chat-server/utils"

	"github.com/tieubaoca/go-chat-server/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageService interface {
	FindMessagesByChatRoomId(requester, chatRoomId string) ([]models.Message, error)
	PaginationMessagesByChatRoomId(requester string, req request.MessagePaginationReq) ([]models.Message, error)
	InsertMessage(message models.Message) (*mongo.InsertOneResult, error)
}

type messageService struct {
	messageRepository       repositories.MessageRepository
	messageStatusRepository repositories.MessageStatusRepository
}

func NewMessageService(
	messageRepository repositories.MessageRepository,
	messageStatusRepository repositories.MessageStatusRepository,
) MessageService {
	return &messageService{
		messageRepository,
		messageStatusRepository,
	}
}

func (s *messageService) FindMessagesByChatRoomId(requester, chatRoomId string) ([]models.Message, error) {
	messages, err := s.messageRepository.FindMessagesByChatRoomId(chatRoomId)
	if err != nil {
		return nil, err
	}
	messageIds := utils.GetMessageIds(messages)
	s.messageStatusRepository.UpdateSeen(requester, messageIds)
	return messages, nil
}

func (s *messageService) InsertMessage(message models.Message) (*mongo.InsertOneResult, error) {
	return s.messageRepository.InsertMessage(message)
}

func (s *messageService) PaginationMessagesByChatRoomId(requester string, req request.MessagePaginationReq) ([]models.Message, error) {
	messages, err := s.messageRepository.PaginationMessagesByChatRoomId(req.ChatRoomId, req.Limit, req.Skip)
	if err != nil {
		return nil, err
	}
	messageIds := utils.GetMessageIds(messages)
	s.messageStatusRepository.UpdateSeen(requester, messageIds)
	return messages, nil
}
