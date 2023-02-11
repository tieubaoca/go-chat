package response

type SaasResponse struct {
	Data       interface{} `json:"data"`
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	Pagination interface{} `json:"pagination"`
}
