package types

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	WebSocketPing = "ping"
	WebSocketPong = "pong"
)

type WebSocketClient struct {
	SaId string
	Conn *websocket.Conn
}

func (wsc *WebSocketClient) Close() {
	wsc.Conn.Close()
}

func (wsc *WebSocketClient) Read() (response.WebSocketEvent, error) {
	var msg response.WebSocketEvent
	err := wsc.Conn.ReadJSON(&msg)

	if err != nil {
		log.Error("error: %v", err)
		wsc.Conn.Close()
		return msg, err
	}
	msg.Sender = wsc.SaId
	msg.Client = wsc.Conn
	return msg, err
}

type OnlineStatusEvent struct {
	SaId     string             `json:"saId"`
	IsActive bool               `json:"isActive"`
	LastSeen primitive.DateTime `json:"lastSeen"`
}
