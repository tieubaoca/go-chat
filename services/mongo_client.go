package services

import (
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func InitDbClient(host string, port string, username string, password string, database string) {
	_dbClient, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/admin", username, password, host, port)))
	if err != nil {
		log.Println(err)
	}
	err = _dbClient.Connect(nil)
	if err != nil {
		log.Println(err)
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
