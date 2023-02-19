package repositories

import (
	"context"
	"time"

	"github.com/tieubaoca/go-chat-server/models"
	"github.com/tieubaoca/go-chat-server/utils/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepository interface {
	FindMessagesByChatRoomId(chatRoomId string) ([]models.Message, error)
	InsertMessage(message models.Message) (*mongo.InsertOneResult, error)
	PaginationMessagesByChatRoomId(chatRoomId string, limit int64, skip int64) ([]models.Message, error)
	UpdateMessageReceivedStatus(messageId []string, saId string) error
	UpdateMessageSeenStatus(messageId []string, saId string) error
	BatchSaIdUpdateMessageReceivedStatus(messageId string, saIds []string) error
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

func (r *messageRepository) UpdateMessageReceivedStatus(messageIds []string, saId string) error {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := r.db.Collection(models.MessageCollection)
	objIds := make([]primitive.ObjectID, len(messageIds))
	for i, v := range messageIds {
		objIds[i], _ = primitive.ObjectIDFromHex(v)
	}
	_, err := coll.UpdateMany(
		context.TODO(),
		bson.M{

			"_id": bson.M{"$in": objIds},
			"receivedBy": bson.M{
				"$not": bson.M{
					"$elemMatch": bson.M{
						"saId": saId,
					},
				},
			},
		},
		bson.M{
			"$push": bson.M{
				"receivedBy": bson.M{
					"saId":       saId,
					"receivedAt": primitive.NewDateTimeFromTime(time.Now()),
				},
			},
		})
	return err
}

func (r *messageRepository) UpdateMessageSeenStatus(messageIds []string, saId string) error {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	if err := r.UpdateMessageReceivedStatus(messageIds, saId); err != nil {
		return err
	}

	objIds := make([]primitive.ObjectID, len(messageIds))
	for i, v := range messageIds {
		objIds[i], _ = primitive.ObjectIDFromHex(v)
	}
	coll := r.db.Collection(models.MessageCollection)

	_, err := coll.UpdateMany(
		context.TODO(),
		bson.M{

			"_id": bson.M{"$in": objIds},
			"seenBy": bson.M{
				"$not": bson.M{
					"$elemMatch": bson.M{
						"saId": saId,
					},
				},
			},
		},
		bson.M{
			"$push": bson.M{
				"seenBy": bson.M{
					"saId":   saId,
					"seenAt": primitive.NewDateTimeFromTime(time.Now()),
				},
			},
		},
	)
	return err
}

func (r *messageRepository) BatchSaIdUpdateMessageReceivedStatus(messageId string, saIds []string) error {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := r.db.Collection(models.MessageCollection)
	objId, err := primitive.ObjectIDFromHex(messageId)
	if err != nil {
		return err
	}
	var message models.Message
	result := coll.FindOne(
		context.TODO(),
		bson.M{
			"_id": objId,
		},
	)
	if result.Err() != nil {
		return result.Err()
	}
	result.Decode(&message)
	updateSaIds := make([]models.Received, 0)
	matchSaIds := make(map[string]bool)
	for _, receive := range message.ReceivedBy {
		matchSaIds[receive.SaId] = true
	}
	for _, saId := range saIds {
		if !matchSaIds[saId] {
			updateSaIds = append(
				updateSaIds,
				models.Received{
					SaId:       saId,
					ReceivedAt: primitive.NewDateTimeFromTime(time.Now()),
				},
			)
		}
	}

	_, err = coll.UpdateOne(
		context.TODO(),
		bson.M{
			"_id": objId,
		},
		bson.M{
			"$push": bson.M{
				"receivedBy": bson.M{
					"$each": updateSaIds,
				},
			},
		},
	)
	return err
}
