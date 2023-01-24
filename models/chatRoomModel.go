package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatroomType string

var (
	ChatroomCollection              = "chatRoom"
	ChatroomTypeDM     ChatroomType = "DM"
	ChatroomTypeGroup  ChatroomType = "GROUP"
)

type ChatRoom struct {
	Id      primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Owner   string             `json:"owner"`
	Type    ChatroomType       `json:"type"`
	Name    string             `json:"name"`
	Members []string           `json:"members"`
}
