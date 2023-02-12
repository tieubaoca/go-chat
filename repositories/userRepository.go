package repositories

import (
	"context"

	"github.com/tieubaoca/go-chat-server/models"
	"github.com/tieubaoca/go-chat-server/utils/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	FindUserStatusInSaIdList(saIds []string) (map[string]models.UserOnlineStatus, error)
	FindUserStatusBySaId(saId string) (models.UserOnlineStatus, error)
	UpdateUserStatus(saId string, isActive bool, lastSeen primitive.DateTime) error
}

type userRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindUserStatusInSaIdList(saIds []string) (map[string]models.UserOnlineStatus, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := r.db.Collection(models.UserOnlineStatusCollection)
	result, err := coll.Find(context.TODO(), bson.D{{"saId", bson.D{{"$in", saIds}}}})
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	var users []models.UserOnlineStatus
	if err = result.All(context.TODO(), &users); err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	mapUser := make(map[string]models.UserOnlineStatus)
	for _, v := range users {
		mapUser[v.SaId] = v
	}
	return mapUser, nil
}

func (r *userRepository) FindUserStatusBySaId(saId string) (models.UserOnlineStatus, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := r.db.Collection(models.UserOnlineStatusCollection)
	result := coll.FindOne(context.TODO(), bson.D{{"saId", saId}})
	var user models.UserOnlineStatus
	if err := result.Decode(&user); err != nil {
		log.ErrorLogger.Println(err)
	}
	return user, nil
}

func (r *userRepository) UpdateUserStatus(saId string, isActive bool, lastSeen primitive.DateTime) error {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := r.db.Collection(models.UserOnlineStatusCollection)
	userStatus, _ := r.FindUserStatusBySaId(saId)
	if userStatus.SaId == "" {
		_, err := coll.InsertOne(context.TODO(), bson.M{
			"saId":     saId,
			"isActive": isActive,
			"lastSeen": lastSeen,
		})
		if err != nil {
			log.ErrorLogger.Println(err)
			return err
		}
		return nil
	}
	_, err := coll.UpdateOne(context.TODO(), bson.D{{"saId", saId}}, bson.D{
		{"$set", bson.D{
			{"isActive", isActive},
			{"lastSeen", lastSeen},
		}},
	})
	if err != nil {
		log.ErrorLogger.Println(err)
		return err
	}
	return nil
}
