package request

type MessagePaginationReq struct {
	ChatRoomId string `json:"chatRoomId"`
	Limit      int64  `json:"limit"`
	Skip       int64  `json:"skip"`
}

type PaginationOnlineFriendReq struct {
	Page int64 `json:"page"`
	Size int64 `json:"size"`
}
