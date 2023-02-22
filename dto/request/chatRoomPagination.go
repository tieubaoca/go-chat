package request

type ChatRoomPagination struct {
	SaId  string `json:"saId" bson:"saId"`
	Skip  int    `json:"skip" bson:"skip"`
	Limit int    `json:"limit" bson:"limit"`
}
