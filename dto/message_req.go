package dto

type MessagePaginationReq struct {
	ChatRoomId string `json:"chatRoomId"`
	Limit      int64  `json:"limit"`
	Skip       int64  `json:"skip"`
}
