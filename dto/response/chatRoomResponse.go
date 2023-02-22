package response

import (
	"github.com/tieubaoca/go-chat-server/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatRoomResponse struct {
	Id          primitive.ObjectID  `bson:"_id" json:"id,omitempty"`
	Owner       string              `json:"owner" bson:"owner"`
	Type        models.ChatRoomType `json:"type" bson:"type"`
	Name        string              `json:"name" bson:"name"`
	Members     []models.Citizen    `json:"members" bson:"members"`
	LastMessage primitive.DateTime  `json:"lastMessage" bson:"lastMessage"`
	IsBlocked   bool                `json:"isBlocked" bson:"isBlocked"`
}
