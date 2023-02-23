package request

type SwitchCitizenReq struct {
	UserId           string `json:"userId" binding:"required"`
	CurrentCitizenId string `json:"currentCitizenId" binding:"required"`
	NewCitizenId     string `json:"newCitizenId" binding:"required"`
}
