package types

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tieubaoca/go-chat-server/dto/response"
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

func (wsc *WebSocketClient) Read() (response.WebSocketResponse, error) {
	var msg response.WebSocketResponse
	err := wsc.Conn.ReadJSON(&msg)

	if err != nil {
		log.Error("error: %v", err)
		wsc.Conn.Close()
		return msg, err
	}
	msg.Sender = wsc.SaId
	return msg, err
}
