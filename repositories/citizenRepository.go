package repositories

import (
	"context"

	"github.com/tieubaoca/go-chat-server/db"
	"github.com/tieubaoca/go-chat-server/models"
	"go.mongodb.org/mongo-driver/bson"
)

type CitizenRepository interface {
	FindCitizenInList(saIds []string) ([]models.Citizen, error)
	FindCitizenBySaId(saId string) (*models.Citizen, error)
}

type citizenRepository struct {
	saasDb *db.SaasMongoDb
}

func NewCitizenRepository(saasDb *db.SaasMongoDb) *citizenRepository {
	return &citizenRepository{
		saasDb: saasDb,
	}
}

func (r *citizenRepository) FindCitizenInList(saIds []string) ([]models.Citizen, error) {
	filter := bson.M{
		"saId": bson.M{
			"$in": saIds,
		},
	}
	coll := r.saasDb.Collection(models.CitizenCollection)
	cur, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var citizens []models.Citizen
	err = cur.All(context.Background(), &citizens)
	if err != nil {
		return nil, err
	}
	return citizens, err
}

func (r *citizenRepository) FindCitizenBySaId(saId string) (*models.Citizen, error) {
	coll := r.saasDb.Collection(models.CitizenCollection)
	filter := bson.M{
		"saId": saId,
	}
	var citizen models.Citizen
	err := coll.FindOne(context.TODO(), filter).Decode(&citizen)
	return &citizen, err
}
