package services

import (
	"net/http"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"

	"github.com/gorilla/websocket"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mapping from session id to websocket client
var wsClients map[string]map[string]*types.WebSocketClient
var upgrader websocket.Upgrader
var broadcast chan response.WebSocketResponse

// / InitWebSocket initializes the websocket server
func InitWebSocket() {
	// Increase the maximum number of open files
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	// Init websocket client map
	wsClients = make(map[string]map[string]*types.WebSocketClient)
	// Init websocket upgrader
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			log.Error("Websocket error: ", reason)
			response.Res(w, types.StatusError, nil, reason.Error())
		},
	}
	broadcast = make(chan response.WebSocketResponse)

	go handleWebSocketResponse()
}

// / HandleWebSocket handles the websocket connection
func HandleWebSocket(w http.ResponseWriter, r *http.Request, saId string, sessionId string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Error(err)
		}
	}()
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	defer ws.Close()

	if _, ok := wsClients[saId]; !ok {
		wsClients[saId] = make(map[string]*types.WebSocketClient)
	}
	log.Info("Add new client: ", saId)
	client := &types.WebSocketClient{
		SaId: saId,
		Conn: ws,
	}
	UpdateUserStatus(saId, true, primitive.NewDateTimeFromTime(time.Now()))

	wsClients[saId][sessionId] = client
	ws.WriteJSON(response.WebSocketResponse{
		EventType:    "Connected",
		EventPayload: saId,
	})
	for {

		var msg response.WebSocketResponse
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Error("error: %v", err)
			client.Close()
			delete(wsClients[saId], sessionId)
			if len(wsClients[saId]) == 0 {
				UpdateUserStatus(saId, false, primitive.NewDateTimeFromTime(time.Now()))
				delete(wsClients, saId)
			}
			break
		}
		msg.Sender = saId
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleWebSocketResponse() {
	defer func() {
		err := recover()
		if err != nil {
			log.Error(err)
		}
	}()
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		switch msg.EventType {
		case types.WebsocketEventTypeMessage:
			handleMessage(msg)
		}
	}
}

func handleMessage(event response.WebSocketResponse) {
	defer func() {
		err := recover()
		if err != nil {
			log.Error(err)
		}
	}()
	msg := event.EventPayload.(map[string]interface{})
	// Send it out to every client that is currently connected
	chatRoom, err := FindChatroomById(msg["chatroom"].(string))
	if err != nil {
		log.Error(err)
	}
	if !utils.ContainsString(chatRoom.Members, event.Sender) {
		for _, ws := range wsClients[event.Sender] {
			ws.Conn.WriteJSON(response.WebSocketResponse{
				EventType:    types.WebsocketEventTypeError,
				Sender:       "server",
				EventPayload: "You are not a member of this chatroom",
			})
		}
		return
	}
	_, err = InsertMessage(bson.M{
		"chatroom": msg["chatroom"].(string),
		"sender":   event.Sender,
		"content":  msg["content"].(string),
		"createAt": primitive.NewDateTimeFromTime(time.Now()),
	})
	if err != nil {
		log.Error(err)
		return
	}
	for _, member := range chatRoom.Members {
		if wss, ok := wsClients[member]; ok {
			for sessionId, ws := range wss {
				err := ws.Conn.WriteJSON(event)
				if err != nil {
					log.Printf("error: %v", err)
					ws.Close()
					delete(wsClients[member], sessionId)
				}
			}

		}
	}

}

func Logout(saId string, sessionId string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Error(err)
		}
	}()
	if _, ok := wsClients[saId]; ok {
		ws, ok := wsClients[saId][sessionId]
		if ok {
			ws.Close()
			delete(wsClients[saId], sessionId)
		}
	}
}
