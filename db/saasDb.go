package db

import "go.mongodb.org/mongo-driver/mongo"

type SaasMongoDb struct {
	*mongo.Database
}

func NewSaasMongoDb(connectionString string, database string) *SaasMongoDb {
	return &SaasMongoDb{
		NewDbClient(connectionString, database),
	}
}
