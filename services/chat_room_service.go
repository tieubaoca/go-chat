package services

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/tieubaoca/go-chat-server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindChatroomById(id string) (models.Chatroom, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.Error(err)
		}
	}()
	coll := db.Collection(models.ChatroomCollection)
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
			log.Error(err)
		}
	}()
	coll := db.Collection(models.ChatroomCollection)
	cursor, err := coll.Find(context.TODO(), bson.D{{"members", member}})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var results []models.Chatroom
	err = cursor.All(context.TODO(), &results)
	return results, err
}

func FindGroupsByMembers(members []string) ([]models.Chatroom, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.Error(err)
		}
	}()
	coll := db.Collection(models.ChatroomCollection)
	var result []models.Chatroom
	cursor, err := coll.Find(
		context.TODO(),
		bson.D{
			{
				"members",
				bson.D{
					{"$all", members},
				},
			},
			{"type", models.ChatroomTypeGroup},
		},
	)
	if err != nil {
		log.Error(err)
		return result, err
	}
	err = cursor.All(context.TODO(), &result)
	return result, err
}

func FindDMByMembers(members []string) (models.Chatroom, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.Error(err)
		}
	}()
	coll := db.Collection(models.ChatroomCollection)
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

func InsertChatroom(chatRoom models.Chatroom) (*mongo.InsertOneResult, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.Error(err)
		}
	}()
	coll := db.Collection(models.ChatroomCollection)
	return coll.InsertOne(
		context.TODO(),
		bson.M{
			"name":    chatRoom.Name,
			"type":    chatRoom.Type,
			"owner":   chatRoom.Owner,
			"members": chatRoom.Members,
		},
	)
}

func AddMemberToChatroom(chatRoomId string, member string) (*mongo.UpdateResult, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.Error(err)
		}
	}()
	coll := db.Collection(models.ChatroomCollection)
	chatRoom, err := FindChatroomById(chatRoomId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if chatRoom.Type == models.ChatroomTypeDM {
		return nil, errors.New("Cannot add member to DM")
	}

	return coll.UpdateOne(context.TODO(), bson.D{{"_id", chatRoom.Id}}, bson.D{
		{"$addToSet", bson.D{{"members", member}}},
	})
}
