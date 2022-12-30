package services

import (
	"context"

	"github.com/tieubaoca/go-chat-server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindMessagesByChatRoomId(chatRoomId string) ([]models.Message, error) {
	coll := db.Collection("message")
	ojId, err := primitive.ObjectIDFromHex(chatRoomId)
	if err != nil {
		return nil, err
	}
	result, err := coll.Find(context.TODO(), bson.D{{"chatRoom", ojId}})
	if err != nil {
		return nil, err
	}
	var messages []models.Message
	if err = result.All(context.TODO(), &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func InsertMessage(message interface{}) (*mongo.InsertOneResult, error) {
	coll := db.Collection("message")
	return coll.InsertOne(context.TODO(), message)
}

func PaginationMessagesByChatRoomId(chatRoomId string, limit int64, skip int64) ([]models.Message, error) {
	coll := db.Collection("message")
	ojId, err := primitive.ObjectIDFromHex(chatRoomId)
	if err != nil {
		return nil, err
	}
	result, err := coll.Find(context.TODO(), bson.D{{"chatRoom", ojId}}, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	})
	if err != nil {
		return nil, err
	}
	var messages []models.Message
	if err = result.All(context.TODO(), &messages); err != nil {
		return nil, err
	}
	return messages, nil
}
