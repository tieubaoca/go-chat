package services

import (
	"net/http"
	"syscall"
	"time"

	"github.com/tieubaoca/go-chat-server/models"
	"github.com/tieubaoca/go-chat-server/utils/log"

	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"

	"github.com/gorilla/websocket"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mapping from session id to websocket client
var wsFd map[string][]int
var wsClients map[int]types.WebSocketClient
var upgrader websocket.Upgrader
var broadcast chan response.WebSocketResponse
var epoll *types.Epoll

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
	wsFd = make(map[string][]int)

	wsClients = make(map[int]types.WebSocketClient)
	epoll, _ = types.MkEpoll()

	// Init websocket upgrader
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	broadcast = make(chan response.WebSocketResponse)
	go HandleEpoll()
	go handleWebSocketResponse()
}

// / HandleWebSocket handles the websocket connection
func HandleWebSocket(w http.ResponseWriter, r *http.Request, saId string) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.ErrorLogger.Println(err)
		return
	}

	// if _, ok := wsClients[saId]; !ok {
	// 	wsClients[saId] = make(map[string]*types.WebSocketClient)
	// }
	log.InfoLogger.Println("Add new client: ", saId)
	client := &types.WebSocketClient{
		SaId: saId,
		Conn: ws,
	}
	// Add new client to epoll
	fd, err := epoll.Add(client.Conn)
	wsClients[fd] = *client
	if err != nil {
		log.ErrorLogger.Println(err)
		return
	}
	UpdateUserStatus(saId, true, primitive.NewDateTimeFromTime(time.Now()))
	if _, ok := wsFd[saId]; !ok {
		wsFd[saId] = make([]int, 0)
	}
	wsFd[saId] = append(wsFd[saId], fd)
	ws.WriteJSON(response.WebSocketResponse{
		EventType:    "Connected",
		EventPayload: saId,
	})
}

func HandleEpoll() {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	for {
		fds, err := epoll.Wait()
		if err != nil {
			log.ErrorLogger.Println(err)
			continue
		}
		for _, fd := range fds {
			client, ok := wsClients[fd]
			if !ok {
				log.ErrorLogger.Println("Client not found")
				continue
			}
			msg, err := client.Read()
			if err != nil {
				log.ErrorLogger.Println(err)
				epoll.Remove(wsClients[fd].Conn)
				delete(wsClients, fd)
				wsFd[client.SaId] = utils.ArrayIntRemoveElement(wsFd[client.SaId], fd)
				if len(wsFd[client.SaId]) == 0 {
					UpdateUserStatus(client.SaId, false, primitive.NewDateTimeFromTime(time.Now()))
				}
				continue
			}
			broadcast <- msg
		}
	}
}

func handleWebSocketResponse() {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
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
			log.ErrorLogger.Println(err)
		}
	}()
	msg := event.EventPayload.(map[string]interface{})
	// Send it out to every client that is currently connected
	chatRoom, err := FindChatRoomById(msg["chatRoom"].(string))
	if err != nil {
		log.ErrorLogger.Println(err)
	}
	if !utils.ContainsString(chatRoom.Members, event.Sender) {
		for _, fd := range wsFd[event.Sender] {
			ws := wsClients[fd]
			ws.Conn.WriteJSON(response.WebSocketResponse{
				EventType:    types.WebsocketEventTypeError,
				Sender:       "server",
				EventPayload: "You are not a member of this chat room",
			})
		}
		return
	}
	message := models.Message{
		ChatRoom: msg["chatRoom"].(string),
		Sender:   event.Sender,
		Content:  msg["content"].(string),
		CreateAt: primitive.NewDateTimeFromTime(time.Now()),
	}
	r, err := InsertMessage(message)
	message.Id = r.InsertedID.(primitive.ObjectID)
	if err != nil {
		log.ErrorLogger.Println(err)
		return
	}
	for _, member := range chatRoom.Members {
		fds := wsFd[member]
		for _, fd := range fds {
			ws := wsClients[fd]
			ws.Conn.WriteJSON(response.WebSocketResponse{
				EventType:    types.WebsocketEventTypeMessage,
				Sender:       event.Sender,
				EventPayload: message,
			})

		}
	}

}

func Logout(saId string) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	if _, ok := wsFd[saId]; ok {
		for _, fd := range wsFd[saId] {
			wsClients[fd].Conn.Close()
			epoll.Remove(wsClients[fd].Conn)
			delete(wsClients, fd)
		}
	}
}
