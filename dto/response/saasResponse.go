package response

type SaasResponse struct {
	Data       interface{} `json:"data"`
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	Pagination interface{} `json:"pagination"`
}

type SaasUserResponse struct {
	Id            interface{} `json:"id"`
	Username      string      `json:"username"`
	SaId          string      `json:"saId"`
	Email         string      `json:"email"`
	Phone         string      `json:"phone"`
	NationalCode  string      `json:"nationalCode"`
	CallingCode   string      `json:"callingCode"`
	FirstName     string      `json:"firstName"`
	LastName      string      `json:"lastName"`
	CreateAt      string      `json:"createAt"`
	UpdateAt      string      `json:"updateAt"`
	ListTopicName []string    `json:"listTopicName"`
}
