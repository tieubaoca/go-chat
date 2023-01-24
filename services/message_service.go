package services

import (
	"context"

	"github.com/tieubaoca/go-chat-server/utils/log"

	"github.com/tieubaoca/go-chat-server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindMessagesByChatRoomId(chatRoomId string) ([]models.Message, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := db.Collection(models.MessageCollection)

	result, err := coll.Find(context.TODO(), bson.D{{"chatroom", chatRoomId}})
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	var messages []models.Message
	if err = result.All(context.TODO(), &messages); err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	return messages, nil
}

func InsertMessage(message interface{}) (*mongo.InsertOneResult, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := db.Collection(models.MessageCollection)
	return coll.InsertOne(context.TODO(), message)
}

func PaginationMessagesByChatRoomId(chatRoomId string, limit int64, skip int64) ([]models.Message, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := db.Collection(models.MessageCollection)
	ojId, err := primitive.ObjectIDFromHex(chatRoomId)
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	result, err := coll.Find(context.TODO(), bson.D{{"chatRoom", ojId}}, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	})
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	var messages []models.Message
	err = result.All(context.TODO(), &messages)
	return messages, nil
}
