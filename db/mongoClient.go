package db

import (
	"context"
	"time"

	"github.com/tieubaoca/go-chat-server/utils/log"

	"github.com/tieubaoca/go-chat-server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDbClient(connectionString string, database string) *mongo.Database {

	opts := options.Client().ApplyURI(connectionString).SetTimeout(2 * time.Second).SetConnectTimeout(3 * time.Second)
	_dbClient, err := mongo.NewClient(opts)
	if err != nil {
		log.ErrorLogger.Panicln(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = _dbClient.Connect(ctx)
	if err != nil {
		log.ErrorLogger.Panicln(err)
	}

	err = _dbClient.Ping(ctx, nil)
	if err != nil {
		log.ErrorLogger.Panicln(err)
	}
	return _dbClient.Database(database)
}

func InitCollections(db *mongo.Database) {
	db.Collection(models.ChatRoomCollection).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"members", 1}},
		Options: options.Index().SetUnique(true),
	})
	db.Collection(models.UserOnlineStatusCollection).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"saId", 1}},
		Options: options.Index().SetUnique(true),
	})
}
