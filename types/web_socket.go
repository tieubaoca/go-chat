package types

import "github.com/gorilla/websocket"

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
