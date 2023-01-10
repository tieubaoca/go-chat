package services

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func InitDbClient(connectionString string, database string) {

	opts := options.Client().ApplyURI(connectionString).SetTimeout(2 * time.Second).SetConnectTimeout(3 * time.Second)
	_dbClient, err := mongo.NewClient(opts)
	if err != nil {
		log.Panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = _dbClient.Connect(ctx)
	if err != nil {
		log.Panic(err)
	}

	db = _dbClient.Database(database)
}

// GetDBClient returns the database client
func GetDBClient() *mongo.Database {
	return db
}

// SetDBClient sets the database client
func SetDBClient(client *mongo.Database) {
	db = client
}

// CloseDBClient closes the database client
func CloseDBClient() {
	db.Client().Disconnect(nil)
}

// Path: services/mongo_client.go
// Compare this snippet from cmd/init.go:
// /*
