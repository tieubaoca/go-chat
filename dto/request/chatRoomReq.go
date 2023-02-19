package request

type AddMemReq struct {
	ChatRoomId string   `json:"chatRoomId"`
	SaIds      []string `json:"saIds"`
}

func (r *AddMemReq) Validate() bool {
	return len(r.ChatRoomId) > 0
}

type RemoveMemReq struct {
	ChatRoomId string   `json:"chatRoomId"`
	SaIds      []string `json:"saIds"`
}
