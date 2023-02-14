package repositories

import (
	"context"

	"github.com/tieubaoca/go-chat-server/models"
	"github.com/tieubaoca/go-chat-server/utils/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepository interface {
	FindMessagesByChatRoomId(chatRoomId string) ([]models.Message, error)
	InsertMessage(message models.Message) (*mongo.InsertOneResult, error)
	PaginationMessagesByChatRoomId(chatRoomId string, limit int64, skip int64) ([]models.Message, error)
}

type messageRepository struct {
	db *mongo.Database
}

func NewMessageRepository(db *mongo.Database) *messageRepository {
	return &messageRepository{db}
}

func (r *messageRepository) FindMessagesByChatRoomId(chatRoomId string) ([]models.Message, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := r.db.Collection(models.MessageCollection)

	result, err := coll.Find(context.TODO(), bson.D{{"chatRoom", chatRoomId}})
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

func (r *messageRepository) InsertMessage(message models.Message) (*mongo.InsertOneResult, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()

	coll := r.db.Collection(models.MessageCollection)
	return coll.InsertOne(context.TODO(), bson.M{
		"chatRoom": message.ChatRoom,
		"sender":   message.Sender,
		"content":  message.Content,
		"createAt": message.CreateAt,
	})
}

func (r *messageRepository) PaginationMessagesByChatRoomId(chatRoomId string, limit int64, skip int64) ([]models.Message, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := r.db.Collection(models.MessageCollection)

	opts := options.Find().SetSkip(skip).SetLimit(limit)
	result, err := coll.Find(context.TODO(), bson.M{
		"chatRoom": chatRoomId,
	}, opts)
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	var messages []models.Message
	err = result.All(context.TODO(), &messages)
	return messages, nil
}
