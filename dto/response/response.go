package response

import "github.com/gorilla/websocket"

type ResponseData struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type WebSocketEvent struct {
	Sender       string          `json:"sender"`
	Client       *websocket.Conn `json:"-"`
	EventType    string          `json:"eventType"`
	EventPayload interface{}     `json:"eventPayload"`
}
