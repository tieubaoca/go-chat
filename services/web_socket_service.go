package services

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tieubaoca/go-chat-server/models"
	"go.mongodb.org/mongo-driver/bson"
)

var wsClients map[string][]*websocket.Conn
var upgrader websocket.Upgrader
var broadcast chan models.Message

func InitWebSocket() {
	wsClients = make(map[string][]*websocket.Conn)
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	broadcast = make(chan models.Message)
	go handleMessage()
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request, username string) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()
	// wsClients[username] = ws
	if _, ok := wsClients[username]; !ok {
		wsClients[username] = []*websocket.Conn{}
	}
	log.Println("Add new client: ", username)
	wsClients[username] = append(wsClients[username], ws)
	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			for i, c := range wsClients[username] {
				if c == ws {
					wsClients[username] = append(wsClients[username][:i], wsClients[username][i+1:]...)
					break
				}
			}
			break
		}
		msg.Sender = username
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessage() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		chatRoom, err := FindChatroomById(msg.Chatroom)
		if err != nil {
			log.Println(err)
		}
		_, err = InsertMessage(bson.M{
			"chat_room": msg.Chatroom,
			"sender":    msg.Sender,
			"content":   msg.Content,
		})
		if err != nil {
			log.Println(err)
			continue
		}
		for _, member := range chatRoom.Members {
			if wss, ok := wsClients[member]; ok {
				for _, ws := range wss {
					err := ws.WriteJSON(msg)
					if err != nil {
						log.Printf("error: %v", err)
						ws.Close()
						delete(wsClients, member)
					}
				}

			}
		}
	}
}
