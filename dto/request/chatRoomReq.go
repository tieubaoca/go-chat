package request

type AddMemReq struct {
	ChatRoomId string   `json:"chatRoomId"`
	SaIds      []string `json:"saIds"`
}

type RemoveMemReq struct {
	ChatRoomId string   `json:"chatRoomId"`
	SaIds      []string `json:"saIds"`
}
