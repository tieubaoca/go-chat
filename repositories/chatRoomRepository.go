package repositories

import (
	"context"
	"errors"

	"github.com/tieubaoca/go-chat-server/models"
	"github.com/tieubaoca/go-chat-server/utils/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRoomRepository interface {
	FindChatRoomById(chatRoomId string) (*models.ChatRoom, error)
	FindChatRoomBySaId(saId string) ([]models.ChatRoom, error)
	PaginationChatRoomBySaId(saId string, skip, limit int) ([]models.ChatRoom, error)
	InsertChatRoom(chatRoom models.ChatRoom) (*mongo.InsertOneResult, error)
	AddMembersToChatRoom(chatRoomId string, members []string) (*mongo.UpdateResult, error)
	RemoveMembersFromChatRoom(chatRoomId string, members []string) (*mongo.UpdateResult, error)
	FindDMByMembers(members []string) (*models.ChatRoom, error)
	FindGroupChatByMembers(members []string) ([]models.ChatRoom, error)
	TransferOwner(chatRoomId, newOwner string) error
	UpdateChatRoomLastMessage(message models.Message) error
}

type chatRoomRepository struct {
	db *mongo.Database
}

func NewChatRoomRepository(db *mongo.Database) *chatRoomRepository {
	return &chatRoomRepository{db}
}

func (r *chatRoomRepository) FindChatRoomById(chatRoomId string) (*models.ChatRoom, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := r.db.Collection(models.ChatRoomCollection)
	var result models.ChatRoom
	obId, err := primitive.ObjectIDFromHex(chatRoomId)
	if err != nil {
		return &result, err
	}
	err = coll.FindOne(context.TODO(), bson.D{{"_id", obId}}).Decode(&result)
	return &result, err
}
func (r *chatRoomRepository) FindChatRoomBySaId(saId string) ([]models.ChatRoom, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := r.db.Collection(models.ChatRoomCollection)
	sortOption := options.Find().SetSort(bson.D{{"lastMessage.createAt", -1}})
	cursor, err := coll.Find(context.TODO(), bson.D{{"members", saId}}, sortOption)
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}

	var results []models.ChatRoom
	err = cursor.All(context.TODO(), &results)
	return results, err
}

func (r *chatRoomRepository) InsertChatRoom(chatRoom models.ChatRoom) (*mongo.InsertOneResult, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := r.db.Collection(models.ChatRoomCollection)
	return coll.InsertOne(
		context.TODO(),
		bson.M{
			"name":      chatRoom.Name,
			"type":      chatRoom.Type,
			"owner":     chatRoom.Owner,
			"isBlocked": false,
			"members":   chatRoom.Members,
		},
	)
}

func (r *chatRoomRepository) AddMembersToChatRoom(chatRoomId string, members []string) (*mongo.UpdateResult, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := r.db.Collection(models.ChatRoomCollection)
	chatRoom, err := r.FindChatRoomById(chatRoomId)

	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	if chatRoom.Type == models.ChatRoomTypeDM {
		return nil, errors.New("Cannot add member to DM")
	}

	return coll.UpdateOne(context.TODO(), bson.M{"_id": chatRoom.Id}, bson.M{
		"$addToSet": bson.M{"members": bson.M{"$each": members}},
	})
}

func (r *chatRoomRepository) RemoveMembersFromChatRoom(chatRoomId string, members []string) (*mongo.UpdateResult, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := r.db.Collection(models.ChatRoomCollection)
	chatRoom, err := r.FindChatRoomById(chatRoomId)
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	if chatRoom.Type == models.ChatRoomTypeDM {
		return nil, errors.New("Cannot remove member from DM")
	}
	return coll.UpdateOne(
		context.TODO(),
		bson.M{"_id": chatRoom.Id},
		bson.M{
			"$pull": bson.M{"members": bson.M{"$in": members}},
		},
	)
}

func (r *chatRoomRepository) FindDMByMembers(saIds []string) (*models.ChatRoom, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()

	coll := r.db.Collection(models.ChatRoomCollection)
	var result models.ChatRoom
	err := coll.FindOne(
		context.TODO(),
		bson.D{
			{
				"members",
				bson.D{
					{"$all", saIds},
				},
			},
			{"type", models.ChatRoomTypeDM},
		},
	).Decode(&result)
	return &result, err
}

func (r *chatRoomRepository) FindGroupChatByMembers(members []string) ([]models.ChatRoom, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := r.db.Collection(models.ChatRoomCollection)
	sortOption := options.Find().SetSort(bson.D{{"lastMessage.createAt", -1}})
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
		sortOption,
	)
	if err != nil {
		log.ErrorLogger.Println(err)
		return result, err
	}
	err = cursor.All(context.TODO(), &result)
	return result, err
}

func (r *chatRoomRepository) TransferOwner(chatRoomId, newOwner string) error {

	coll := r.db.Collection(models.ChatRoomCollection)
	objId, err := primitive.ObjectIDFromHex(chatRoomId)
	if err != nil {
		return err
	}
	result := coll.FindOneAndUpdate(
		context.TODO(),
		bson.M{
			"_id": objId,
		},
		bson.M{
			"$set": bson.M{
				"owner": newOwner,
			},
		},
	)

	return result.Err()
}

func (r *chatRoomRepository) UpdateChatRoomLastMessage(message models.Message) error {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := r.db.Collection(models.ChatRoomCollection)
	objId, err := primitive.ObjectIDFromHex(message.ChatRoom)
	if err != nil {
		return err
	}
	return coll.FindOneAndUpdate(
		context.TODO(),
		bson.M{
			"_id": objId,
		},
		bson.M{
			"$set": bson.M{
				"lastMessage": message,
			},
		},
	).Err()
}

func (r *chatRoomRepository) PaginationChatRoomBySaId(saId string, skip, limit int) ([]models.ChatRoom, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()

	coll := r.db.Collection(models.ChatRoomCollection)
	var result []models.ChatRoom
	option := options.Find().SetSort(bson.M{"lastMessage.createAt": -1}).SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := coll.Find(
		context.TODO(),
		bson.M{
			"members": bson.M{
				"$elemMatch": bson.M{
					"$eq": saId,
				},
			},
		},
		option,
	)
	if err != nil {
		log.ErrorLogger.Println(err)
		return result, err
	}
	err = cursor.All(context.TODO(), &result)
	return result, err

}
