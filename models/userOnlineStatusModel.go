package models

import "go.mongodb.org/mongo-driver/bson/primitive"

var UserOnlineStatusCollection = "userOnlineStatus"

type UserOnlineStatus struct {
	Id       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	SaId     string             `json:"saId" bson:"saId"`
	IsActive bool               `json:"isActive" bson:"isActive"`
	LastSeen primitive.DateTime `json:"lastSeen" bson:"lastSeen"`
}
