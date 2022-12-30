package dto

type GetAccessTokenReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
