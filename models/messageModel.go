package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var MessageCollection = "message"

type Message struct {
	Id         primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	ChatRoom   string             `json:"chatRoom" bson:"chatRoom"`
	Sender     string             `json:"sender" bson:"sender"`
	Content    string             `json:"content" bson:"content"`
	CreateAt   primitive.DateTime `json:"createAt" bson:"createAt"`
	ReceivedBy []Received         `json:"receivedBy" bson:"receivedBy"`
	SeenBy     []Seen             `json:"seenBy" bson:"seenBy"`
}

type Seen struct {
	SaId   string             `json:"saId" bson:"saId"`
	SeenAt primitive.DateTime `json:"seenAt" bson:"seenAt"`
}

type Received struct {
	SaId       string             `json:"saId" bson:"saId"`
	ReceivedAt primitive.DateTime `json:"receivedAt" bson:"receivedAt"`
}
