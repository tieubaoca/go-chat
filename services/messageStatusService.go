package services

import (
	"github.com/tieubaoca/go-chat-server/models"
	"github.com/tieubaoca/go-chat-server/repositories"
)

type MessageStatusService interface {
	FindMessageStatusByMessageId(messageId string) (*models.MessageStatus, error)
	FindMessageStatusByMessageIds(messageIds []string) ([]models.MessageStatus, error)
}

type messageStatusService struct {
	messageStatusRepository repositories.MessageStatusRepository
}

func NewMessageStatusService(
	messageStatusRepository repositories.MessageStatusRepository,
) *messageStatusService {
	return &messageStatusService{
		messageStatusRepository: messageStatusRepository,
	}
}

func (s *messageStatusService) FindMessageStatusByMessageId(messageId string) (*models.MessageStatus, error) {
	return s.messageStatusRepository.FindMessageStatusByMessageId(messageId)
}

func (s *messageStatusService) FindMessageStatusByMessageIds(messageIds []string) ([]models.MessageStatus, error) {
	return s.messageStatusRepository.FindMessageStatusInListMessageId(messageIds)
}
