package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatroomType string

var (
	ChatroomTypeDM    ChatroomType = "DM"
	ChatroomTypeGroup ChatroomType = "GROUP"
)

type Chatroom struct {
	Id      primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Owner   string             `json:"owner"`
	Type    ChatroomType       `json:"type"`
	Name    string             `json:"name"`
	Members []string           `json:"members"`
}
