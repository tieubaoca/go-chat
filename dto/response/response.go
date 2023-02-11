package response

type ResponseData struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type WebSocketResponse struct {
	Sender       string      `json:"sender"`
	EventType    string      `json:"eventType"`
	EventPayload interface{} `json:"eventPayload"`
}
