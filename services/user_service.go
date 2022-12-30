package services

import (
	"context"

	"github.com/tieubaoca/go-chat-server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindUserByUsername(username string) (models.User, error) {
	coll := GetDBClient().Collection("user")
	var result models.User
	err := coll.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&result)
	return result, err
}

func FindUserById(id string) (models.User, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, err
	}
	coll := db.Collection("user")
	var result models.User
	err = coll.FindOne(context.TODO(), bson.D{{"_id", objId}}).Decode(&result)
	return result, err
}

func InsertUser(user interface{}) (*mongo.InsertOneResult, error) {
	coll := db.Collection("user")
	return coll.InsertOne(context.TODO(), user)

}

func FindOnlineUsers() ([]models.User, error) {
	users := make([]models.User, 0)
	for username, cs := range wsClients {
		if len(cs) > 0 {
			users = append(users, models.User{Username: username})
		}
	}
	return users, nil
}
