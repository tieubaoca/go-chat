package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var MessageCollection = "message"

type Message struct {
	Id       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Chatroom string             `json:"chatroom"`
	Sender   string             `json:"sender"`
	Content  string             `json:"content"`
	CreateAt primitive.DateTime `json:"createAt"`
}
