package services

import (
	"github.com/tieubaoca/go-chat-server/repositories"
	"github.com/tieubaoca/go-chat-server/utils/log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	FindUserStatusInUserList(usersInfo []interface{}) ([]interface{}, error)
	FindUserStatusBySaId(userInfo map[string]interface{}) (map[string]interface{}, error)
	UpdateUserStatus(saId string, isActive bool, lastSeen primitive.DateTime) error
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

// input: friends info response from saas api
func (s *userService) FindUserStatusInUserList(usersInfo []interface{}) ([]interface{}, error) {
	saIds := make([]string, 0)
	for _, v := range usersInfo {
		saIds = append(saIds, v.(map[string]interface{})["saId"].(string))
	}
	mapFriendStatus, err := s.userRepository.FindUserStatusInSaIdList(saIds)
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	for _, friend := range usersInfo {
		if friendStatus, ok := mapFriendStatus[friend.(map[string]interface{})["saId"].(string)]; ok {
			friend.(map[string]interface{})["lastSeen"] = friendStatus.LastSeen.Time()
			friend.(map[string]interface{})["isActive"] = friendStatus.IsActive
		} else {
			friend.(map[string]interface{})["lastSeen"] = ""
			friend.(map[string]interface{})["isActive"] = false
		}
	}
	return usersInfo, nil
}

func (s *userService) FindUserStatusBySaId(userInfo map[string]interface{}) (map[string]interface{}, error) {
	saId := userInfo["saId"].(string)
	friendStatus, err := s.userRepository.FindUserStatusBySaId(saId)
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	userInfo["lastSeen"] = friendStatus.LastSeen.Time()
	userInfo["isActive"] = friendStatus.IsActive
	return userInfo, nil
}

func (s *userService) UpdateUserStatus(saId string, isActive bool, lastSeen primitive.DateTime) error {
	return s.userRepository.UpdateUserStatus(saId, isActive, lastSeen)
}
