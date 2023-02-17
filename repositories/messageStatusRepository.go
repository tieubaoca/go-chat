package repositories

import (
	"context"
	"time"

	"github.com/tieubaoca/go-chat-server/models"
	"github.com/tieubaoca/go-chat-server/utils/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageStatusRepository interface {
	FindMessageStatusByMessageId(messageId string) (*models.MessageStatus, error)
	FindMessageStatusInListMessageId(messageIds []string) ([]models.MessageStatus, error)
	UpdateSeen(saId string, messageIds []string) error
	UpdateReceived(saId string, messageIds []string) error
	UpdateReceivedBatchSaIds(saIds []string, messageId string) error
}

type messageStatusRepository struct {
	db *mongo.Database
}

func NewMessageStatusRepository(db *mongo.Database) *messageStatusRepository {
	return &messageStatusRepository{db: db}
}

func (r *messageStatusRepository) FindMessageStatusByMessageId(messageId string) (*models.MessageStatus, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := r.db.Collection(models.MessageStatusCollection)
	result := coll.FindOne(context.TODO(), bson.D{{"messageId", messageId}})
	if result.Err() != nil {
		log.ErrorLogger.Println(result.Err())
		return nil, result.Err()
	}
	var messageStatus models.MessageStatus
	if err := result.Decode(&messageStatus); err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	return &messageStatus, nil
}

func (r *messageStatusRepository) FindMessageStatusInListMessageId(messageIds []string) ([]models.MessageStatus, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := r.db.Collection(models.MessageStatusCollection)
	result, err := coll.Find(context.TODO(), bson.D{{"messageId", bson.D{{"$in", messageIds}}}})
	if err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	var messageStatuses []models.MessageStatus
	if err = result.All(context.TODO(), &messageStatuses); err != nil {
		log.ErrorLogger.Println(err)
		return nil, err
	}
	return messageStatuses, nil
}

func (r *messageStatusRepository) UpdateSeen(saId string, messageIds []string) error {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	r.UpdateReceived(saId, messageIds)
	coll := r.db.Collection(models.MessageStatusCollection)
	_, err := coll.UpdateMany(
		context.TODO(),
		bson.M{

			"messageId": bson.M{"$in": messageIds},
			"seenBy": bson.M{
				"saId": bson.M{
					"$nin": []string{saId},
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

func (r *messageStatusRepository) UpdateReceived(saId string, messageIds []string) error {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := r.db.Collection(models.MessageStatusCollection)
	_, err := coll.UpdateMany(
		context.TODO(),
		bson.M{

			"messageId": bson.M{"$in": messageIds},
			"receivedBy": bson.M{
				"saId": bson.M{
					"$nin": []string{saId},
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
		},
	)
	return err
}

func (r *messageStatusRepository) UpdateReceivedBatchSaIds(saIds []string, messageId string) error {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	coll := r.db.Collection(models.MessageStatusCollection)
	result := coll.FindOne(context.TODO(), bson.M{
		"messageId": messageId,
	})
	if result.Err() != nil {
		log.ErrorLogger.Println(result.Err())
		return result.Err()
	}
	var messageStatus models.MessageStatus
	if err := result.Decode(&messageStatus); err != nil {
		log.ErrorLogger.Println(err)
		return err
	}
	saIdUpdate := make([]models.Received, 0)
	for _, saId := range saIds {
		if messageStatus.IsReceivedBy(saId) {
			continue
		}
		saIdUpdate = append(saIdUpdate, models.Received{
			SaId:       saId,
			ReceivedAt: primitive.NewDateTimeFromTime(time.Now()),
		})
	}
	_, err := coll.UpdateOne(
		context.TODO(),
		bson.M{
			"messageId": messageId,
		},
		bson.M{
			"$push": bson.M{
				"receivedBy": saIdUpdate,
			},
		},
	)
	return err

}
