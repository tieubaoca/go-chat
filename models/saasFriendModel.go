package models

const SaasFriendCollection = "friend"

type SaasFriend struct {
	Id         string `json:"id" bson:"_id"`
	SaId       string `json:"saId" bson:"saId"`
	SaIdFriend string `json:"saIdFriend" bson:"saIdFriend"`
	CreateAt   int64  `json:"createAt" bson:"createAt"`
}
