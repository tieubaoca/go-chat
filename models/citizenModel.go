package models

const CitizenCollection = "citizen"

type Citizen struct {
	Id     string `json:"id" bson:"_id"`
	SaId   string `json:"saId" bson:"saId"`
	UserId string `json:"userId" bson:"userId"`
	Name   string `json:"name" bson:"name"`
}
