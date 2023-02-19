package request

type GetAccessTokenReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *GetAccessTokenReq) Validate() bool {
	return len(r.Username) > 0 && len(r.Password) > 0
}
