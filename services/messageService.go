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
	UpdateMessageReceivedStatus(messageId []string, saId string) error
	UpdateMessageSeenStatus(messageId []string, saId string) error
}

type messageService struct {
	messageRepository repositories.MessageRepository
}

func NewMessageService(
	messageRepository repositories.MessageRepository,
) MessageService {
	return &messageService{
		messageRepository,
	}
}

func (s *messageService) FindMessagesByChatRoomId(requester, chatRoomId string) ([]models.Message, error) {
	messages, err := s.messageRepository.FindMessagesByChatRoomId(chatRoomId)
	if err != nil {
		return nil, err
	}
	messageIds := utils.GetMessageIds(messages)
	s.UpdateMessageReceivedStatus(messageIds, requester)
	s.UpdateMessageSeenStatus(messageIds, requester)
	return messages, nil
}

func (s *messageService) InsertMessage(message models.Message) (*mongo.InsertOneResult, error) {

	result, err := s.messageRepository.InsertMessage(message)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *messageService) PaginationMessagesByChatRoomId(requester string, req request.MessagePaginationReq) ([]models.Message, error) {
	messages, err := s.messageRepository.PaginationMessagesByChatRoomId(req.ChatRoomId, req.Limit, req.Skip)
	if err != nil {
		return nil, err
	}
	messageIds := utils.GetMessageIds(messages)
	s.UpdateMessageReceivedStatus(messageIds, requester)
	s.UpdateMessageSeenStatus(messageIds, requester)
	return messages, nil
}

func (s *messageService) UpdateMessageReceivedStatus(messageId []string, saId string) error {
	return s.messageRepository.UpdateMessageReceivedStatus(
		messageId,
		saId,
	)
}

func (s *messageService) UpdateMessageSeenStatus(messageId []string, saId string) error {
	return s.messageRepository.UpdateMessageSeenStatus(messageId, saId)
}

// func (s *messageService) updateSeen(requester string, messageIds []string) error {
// 	notExist := s.messageStatusRepository.GetNotExist(messageIds)
// 	log.InfoLogger.Println(notExist)
// 	if len(notExist) > 0 {
// 		s.messageStatusRepository.InsertManyMessageStatus(
// 			notExist,
// 		)
// 	}
// 	return s.messageStatusRepository.UpdateSeen(requester, messageIds)
// }
