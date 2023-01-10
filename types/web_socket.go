package types

import "github.com/gorilla/websocket"

type WebSocketClient struct {
	SaId string
	Conn *websocket.Conn
}

func (wsc *WebSocketClient) Close() {
	wsc.Conn.Close()
}
