package services

import (
	"context"
	"errors"

	"github.com/tieubaoca/go-chat-server/utils/log"

	"github.com/tieubaoca/go-chat-server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindChatRoomById(id string) (models.ChatRoom, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := db.Collection(models.ChatRoomCollection)
	var result models.ChatRoom
	obId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}
	err = coll.FindOne(context.TODO(), bson.D{{"_id", obId}}).Decode(&result)
	return result, err
}

func FindChatRoomsByMember(member string) ([]models.ChatRoom, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := db.Collection(models.ChatRoomCollection)
	cursor, err := coll.Find(context.TODO(), bson.D{{"members", member}})
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}

	var results []models.ChatRoom
	err = cursor.All(context.TODO(), &results)
	return results, err
}

func FindGroupsByMembers(members []string) ([]models.ChatRoom, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := db.Collection(models.ChatRoomCollection)
	var result []models.ChatRoom
	cursor, err := coll.Find(
		context.TODO(),
		bson.D{
			{
				"members",
				bson.D{
					{"$all", members},
				},
			},
			{"type", models.ChatRoomTypeGroup},
		},
	)
	if err != nil {
		log.ErrorLogger.Println(err)
		return result, err
	}
	err = cursor.All(context.TODO(), &result)
	return result, err
}

func FindDMByMembers(members []string) (models.ChatRoom, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := db.Collection(models.ChatRoomCollection)
	var result models.ChatRoom
	err := coll.FindOne(
		context.TODO(),
		bson.D{
			{
				"members",
				bson.D{
					{"$all", members},
				},
			},
			{"type", models.ChatRoomTypeDM},
		},
	).Decode(&result)
	return result, err
}

func InsertChatRoom(chatRoom models.ChatRoom) (*mongo.InsertOneResult, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := db.Collection(models.ChatRoomCollection)
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

func AddMemberToChatRoom(chatRoomId string, member string) (*mongo.UpdateResult, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := db.Collection(models.ChatRoomCollection)
	chatRoom, err := FindChatRoomById(chatRoomId)
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	if chatRoom.Type == models.ChatRoomTypeDM {
		return nil, errors.New("Cannot add member to DM")
	}

	return coll.UpdateOne(context.TODO(), bson.D{{"_id", chatRoom.Id}}, bson.D{
		{"$addToSet", bson.D{{"members", member}}},
	})
}
