package services

import (
	"context"

	"github.com/tieubaoca/go-chat-server/utils/log"

	"github.com/tieubaoca/go-chat-server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindUserStatusInSaIdList(saIds []string) map[string]models.UserOnlineStatus {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := db.Collection(models.UserOnlineStatusCollection)
	result, err := coll.Find(context.TODO(), bson.D{{"saId", bson.D{{"$in", saIds}}}})
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil
	}
	var users []models.UserOnlineStatus
	if err = result.All(context.TODO(), &users); err != nil {
		log.ErrorLogger.Println(err)
		return nil
	}
	mapUser := make(map[string]models.UserOnlineStatus)
	for _, v := range users {
		mapUser[v.SaId] = v
	}
	return mapUser

}

func FindUserStatusBySaId(saId string) models.UserOnlineStatus {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := db.Collection(models.UserOnlineStatusCollection)
	result := coll.FindOne(context.TODO(), bson.D{{"saId", saId}})
	var user models.UserOnlineStatus
	if err := result.Decode(&user); err != nil {
		log.ErrorLogger.Println(err)
	}
	return user
}

func UpdateUserStatus(saId string, isActive bool, lastSeen primitive.DateTime) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := db.Collection(models.UserOnlineStatusCollection)
	if FindUserStatusBySaId(saId).SaId == "" {
		_, err := coll.InsertOne(context.TODO(), bson.M{
			"saId":     saId,
			"isActive": isActive,
			"lastSeen": lastSeen,
		})
		if err != nil {
			log.ErrorLogger.Println(err)
		}
		return
	}
	_, err := coll.UpdateOne(context.TODO(), bson.D{{"saId", saId}}, bson.D{
		{"$set", bson.D{
			{"isActive", isActive},
			{"lastSeen", lastSeen},
		}},
	})
	if err != nil {
		log.ErrorLogger.Println(err)
	}

}
