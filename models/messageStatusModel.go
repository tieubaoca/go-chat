package models

import "go.mongodb.org/mongo-driver/bson/primitive"

var MessageStatusCollection = "messageStatus"

type MessageStatus struct {
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	MessageId  string             `json:"messageId" bson:"messageId"`
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

func (m *MessageStatus) IsSeenBy(saId string) bool {
	for _, seen := range m.SeenBy {
		if seen.SaId == saId {
			return true
		}
	}
	return false
}

func (m *MessageStatus) IsReceivedBy(saId string) bool {
	for _, received := range m.ReceivedBy {
		if received.SaId == saId {
			return true
		}
	}
	return false
}
