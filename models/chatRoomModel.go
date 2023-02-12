package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatRoomType string

var (
	ChatRoomCollection              = "chatRoom"
	ChatRoomTypeDM     ChatRoomType = "DM"
	ChatRoomTypeGroup  ChatRoomType = "GROUP"
)

type ChatRoom struct {
	Id        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Owner     string             `json:"owner" bson:"owner"`
	Type      ChatRoomType       `json:"type" bson:"type"`
	Name      string             `json:"name" bson:"name"`
	Members   []string           `json:"members" bson:"members"`
	IsBlocked bool               `json:"isBlocked" bson:"isBlocked"`
}
