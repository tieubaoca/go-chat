package response

import (
	"encoding/json"
	"net/http"
)

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

func Res(w http.ResponseWriter, status string, data interface{}, message string) {
	res := ResponseData{
		Status:  status,
		Data:    data,
		Message: message,
	}
	json.NewEncoder(w).Encode(res)
}
