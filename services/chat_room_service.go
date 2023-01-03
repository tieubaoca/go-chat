package services

import (
	"context"
	"log"

	"github.com/tieubaoca/go-chat-server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindChatroomById(id string) (models.Chatroom, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()
	coll := db.Collection("chat_room")
	var result models.Chatroom
	obId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}
	err = coll.FindOne(context.TODO(), bson.D{{"_id", obId}}).Decode(&result)
	return result, err
}

func FindChatroomsByMember(member string) ([]models.Chatroom, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()
	coll := db.Collection("chat_room")
	cursor, err := coll.Find(context.TODO(), bson.D{{"members", member}})
	if err != nil {
		return nil, err
	}

	var results []models.Chatroom
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return results, err
}

func FindDMByMembers(members []string) (models.Chatroom, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()
	coll := db.Collection("chat_room")
	var result models.Chatroom
	err := coll.FindOne(
		context.TODO(),
		bson.D{
			{
				"members",
				bson.D{
					{"$all", members},
				},
			},
			{"type", models.ChatroomTypeDM},
		},
	).Decode(&result)
	return result, err
}

func InsertChatroom(chatRoom interface{}) (*mongo.InsertOneResult, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()
	coll := db.Collection("chat_room")
	return coll.InsertOne(context.TODO(), chatRoom)
}

func AddMemberToChatroom(chatRoomId string, member string) (*mongo.UpdateResult, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()
	coll := db.Collection("chat_room")
	obId, err := primitive.ObjectIDFromHex(chatRoomId)
	if err != nil {
		return nil, err
	}
	return coll.UpdateOne(context.TODO(), bson.D{{"_id", obId}}, bson.D{{"$addToSet", bson.D{{"members", member}}}})
}
