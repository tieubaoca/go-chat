package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var MessageCollection = "message"

type Message struct {
	Id        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	ChatRoom  string             `json:"chatRoom" bson:"chatRoom"`
	Sender    string             `json:"sender" bson:"sender"`
	Content   string             `json:"content" bson:"content"`
	CreateAt  primitive.DateTime `json:"createAt" bson:"createAt"`
	IsBlocked bool               `json:"isBlocked" bson:"isBlocked"`
}
