package repositories

import (
	"context"

	"github.com/tieubaoca/go-chat-server/db"
	"github.com/tieubaoca/go-chat-server/models"
	"go.mongodb.org/mongo-driver/bson"
)

type SaasFriendRepository interface {
	FindAllBySaId(saId string) ([]models.SaasFriend, error)
}

type saasFriendRepository struct {
	db *db.SaasMongoDb
}

func NewSaasFriendRepository(db *db.SaasMongoDb) *saasFriendRepository {
	return &saasFriendRepository{
		db: db,
	}
}

func (r *saasFriendRepository) FindAllBySaId(saId string) ([]models.SaasFriend, error) {
	coll := r.db.Collection("saas_friend")
	var saasFriends []models.SaasFriend
	cur, err := coll.Find(context.Background(), bson.M{"saId": saId})
	if err != nil {
		return nil, err
	}
	if err := cur.All(context.TODO(), &saasFriends); err != nil {
		return nil, err
	}
	return saasFriends, nil
}
