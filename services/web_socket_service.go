package services

import (
	"log"
	"net/http"
	"syscall"

	"github.com/tieubaoca/go-chat-server/saconstant"
	"github.com/tieubaoca/go-chat-server/types"

	"github.com/gorilla/websocket"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"go.mongodb.org/mongo-driver/bson"
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
			log.Println("Websocket error: ", reason)
			response.Res(w, saconstant.StatusError, nil, reason.Error())
		},
	}
	broadcast = make(chan response.WebSocketResponse)

	go handleWebSocketResponse()
}

// / HandleWebSocket handles the websocket connection
func HandleWebSocket(w http.ResponseWriter, r *http.Request, username string, sessionId string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	if _, ok := wsClients[username]; !ok {
		wsClients[username] = make(map[string]*types.WebSocketClient)
	}
	log.Println("Add new client: ", username)
	client := &types.WebSocketClient{
		Username: username,
		Conn:     ws,
	}
	wsClients[username][sessionId] = client
	ws.WriteJSON(response.WebSocketResponse{
		EventType:    "Connected",
		EventPayload: username,
	})
	for {

		var msg response.WebSocketResponse
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(wsClients[username], sessionId)
			break
		}
		msg.Sender = username
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleWebSocketResponse() {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		switch msg.EventType {
		case "Message":
			handleMessage(msg)
		}
	}
}

func handleMessage(event response.WebSocketResponse) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()
	msg := event.EventPayload.(map[string]interface{})
	// Send it out to every client that is currently connected
	chatRoom, err := FindChatroomById(msg["chatroom"].(string))
	if err != nil {
		log.Println(err)
	}
	_, err = InsertMessage(bson.M{
		"chatroom": msg["chatroom"].(string),
		"sender":   event.Sender,
		"content":  msg["content"].(string),
	})
	if err != nil {
		log.Println(err)
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

func Logout(username string, sessionId string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()
	if _, ok := wsClients[username]; ok {
		ws, ok := wsClients[username][sessionId]
		if ok {
			ws.Close()
			delete(wsClients[username], sessionId)
		}
	}
}
