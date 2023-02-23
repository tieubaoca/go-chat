package services

import (
	"encoding/json"
	"net/http"
	"syscall"
	"time"

	"github.com/tieubaoca/go-chat-server/models"
	"github.com/tieubaoca/go-chat-server/repositories"
	"github.com/tieubaoca/go-chat-server/utils/log"

	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils"

	"github.com/gorilla/websocket"
	"github.com/tieubaoca/go-chat-server/dto/request"
	"github.com/tieubaoca/go-chat-server/dto/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebSocketService interface {
	HandleWebSocket(w http.ResponseWriter, r *http.Request, saId string)
	HandleEpoll()
	Logout(saId string)
	SwitchCitizen(switchCitizenReq request.SwitchCitizenReq) error
}

type webSocketService struct {
	chatRoomRepository repositories.ChatRoomRepository
	messageRepository  repositories.MessageRepository
	userRepository     repositories.UserRepository
	epoll              *types.Epoll
	wsFd               map[string][]int
	wsClients          map[int]*types.WebSocketClient
	upgrader           websocket.Upgrader
}

// / InitWebSocket initializes the websocket server
func NewWebSocketService(
	chatRoomRepository repositories.ChatRoomRepository,
	messageRepository repositories.MessageRepository,
	userRepository repositories.UserRepository,

) *webSocketService {
	// Increase the maximum number of open files
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	epoll, err := types.MkEpoll()
	if err != nil {
		panic(err)
	}
	wsService := &webSocketService{
		chatRoomRepository: chatRoomRepository,
		messageRepository:  messageRepository,
		userRepository:     userRepository,
		epoll:              epoll,
		wsFd:               make(map[string][]int),
		wsClients:          make(map[int]*types.WebSocketClient),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	go wsService.HandleEpoll()
	return wsService
}

// / HandleWebSocket handles the websocket connection
func (s *webSocketService) HandleWebSocket(w http.ResponseWriter, r *http.Request, saId string) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.ErrorLogger.Println(err)
		return
	}
	client := &types.WebSocketClient{
		SaId: saId,
		Conn: ws,
	}
	// Add new client to epoll
	fd, err := s.epoll.Add(client.Conn)
	s.wsClients[fd] = client
	if err != nil {
		log.ErrorLogger.Println(err)
		return
	}
	s.userRepository.UpdateUserStatus(saId, true, primitive.NewDateTimeFromTime(time.Now()))
	if _, ok := s.wsFd[saId]; !ok {
		s.wsFd[saId] = make([]int, 0)
	}
	s.wsFd[saId] = append(s.wsFd[saId], fd)
	ws.WriteJSON(response.WebSocketEvent{
		EventType:    "Connected",
		EventPayload: saId,
	})
}

func (s *webSocketService) HandleEpoll() {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	for {
		fds, err := s.epoll.Wait()
		if err != nil {
			continue
		}
		for _, fd := range fds {
			client, ok := s.wsClients[fd]
			if !ok {
				log.ErrorLogger.Println("Client not found")
				continue
			}
			msg, err := client.Read()
			if err != nil {
				log.ErrorLogger.Println(err)
				s.epoll.Remove(s.wsClients[fd].Conn)
				delete(s.wsClients, fd)
				s.wsFd[client.SaId] = utils.ArrayIntRemoveElement(s.wsFd[client.SaId], fd)
				if len(s.wsFd[client.SaId]) == 0 {
					s.userRepository.UpdateUserStatus(client.SaId, false, primitive.NewDateTimeFromTime(time.Now()))
				}
				continue
			}
			switch msg.EventType {
			case types.WebsocketEventTypeMessage:
				go s.handleMessage(msg)
			case types.WebSocketEventTypeTyping:
				go s.handleTypingEvent(msg)
			default:
				go s.handleError(msg)
			}
		}
	}
}

func (s *webSocketService) handleMessage(event response.WebSocketEvent) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	msg := event.EventPayload.(map[string]interface{})
	// Send it out to every client that is currently connected
	chatRoom, err := s.chatRoomRepository.FindChatRoomById(msg["chatRoom"].(string))
	if err != nil {
		log.ErrorLogger.Println(err)
		event.Client.WriteJSON(response.WebSocketEvent{
			EventType:    types.WebsocketEventTypeError,
			Sender:       "server",
			EventPayload: err.Error(),
		})
		return
	}
	if !utils.ContainsString(chatRoom.Members, event.Sender) {
		event.Client.WriteJSON(response.WebSocketEvent{
			EventType:    types.WebsocketEventTypeError,
			Sender:       "server",
			EventPayload: "You are not a member of this chat room",
		})

		return
	}
	message := models.Message{
		ChatRoom:   msg["chatRoom"].(string),
		Sender:     event.Sender,
		Content:    msg["content"].(string),
		CreateAt:   primitive.NewDateTimeFromTime(time.Now()),
		ReceivedBy: make([]models.Received, 0),
		SeenBy:     make([]models.Seen, 0),
	}
	result, err := s.messageRepository.InsertMessage(message)

	message.Id = result.InsertedID.(primitive.ObjectID)
	if err != nil {
		log.ErrorLogger.Println(err)
		return
	}
	err = s.chatRoomRepository.UpdateChatRoomLastMessage(message)
	if err != nil {
		log.ErrorLogger.Println(err)
	}
	citizenEmitted := make([]string, 0)
	for _, member := range chatRoom.Members {
		fds := s.wsFd[member]
		for _, fd := range fds {
			ws := s.wsClients[fd]
			err := ws.Conn.WriteJSON(response.WebSocketEvent{
				EventType:    types.WebsocketEventTypeMessage,
				Sender:       event.Sender,
				EventPayload: message,
			})
			if err != nil {
				log.ErrorLogger.Println(err)
				s.epoll.Remove(s.wsClients[fd].Conn)
				continue
			}
		}
		citizenEmitted = append(citizenEmitted, member)
	}
	s.messageRepository.BatchSaIdUpdateMessageReceivedStatus(message.Id.Hex(), citizenEmitted)
}

func (s *webSocketService) handleTypingEvent(event response.WebSocketEvent) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	var typingEvent request.WebsocketTypingEvent
	jsonString, err := json.Marshal(event.EventPayload)
	if err != nil {
		event.Client.WriteJSON(
			response.WebSocketEvent{
				Sender:       "server",
				EventType:    types.WebsocketEventTypeError,
				EventPayload: err,
			},
		)
		return
	}
	json.Unmarshal(jsonString, &typingEvent)
	chatRoom, err := s.chatRoomRepository.FindChatRoomById(typingEvent.ChatRoomId)
	if err != nil {
		log.ErrorLogger.Println(err)
		event.Client.WriteJSON(response.WebSocketEvent{
			EventType:    types.WebsocketEventTypeError,
			Sender:       "server",
			EventPayload: err.Error(),
		})
		return
	}
	if !utils.ContainsString(chatRoom.Members, event.Sender) {
		event.Client.WriteJSON(response.WebSocketEvent{
			EventType:    types.WebsocketEventTypeError,
			Sender:       "server",
			EventPayload: types.ErrorNotRoomMember,
		})
		return
	}
	for _, member := range chatRoom.Members {
		fds := s.wsFd[member]
		for _, fd := range fds {
			ws := s.wsClients[fd]
			err := ws.Conn.WriteJSON(response.WebSocketEvent{
				EventType:    types.WebSocketEventTypeTyping,
				Sender:       event.Sender,
				EventPayload: typingEvent,
			})
			if err != nil {
				log.ErrorLogger.Println(err)
				s.epoll.Remove(s.wsClients[fd].Conn)
				continue
			}
		}
	}
}

func (s *webSocketService) handleError(event response.WebSocketEvent) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	log.ErrorLogger.Println(event.EventPayload)
	event.Client.WriteJSON(response.WebSocketEvent{
		EventType:    types.WebsocketEventTypeError,
		Sender:       "server",
		EventPayload: "invalid event type",
	})
}

func (s *webSocketService) Logout(saId string) {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()
	if _, ok := s.wsFd[saId]; ok {
		for _, fd := range s.wsFd[saId] {
			s.wsClients[fd].Conn.Close()
			s.epoll.Remove(s.wsClients[fd].Conn)
			delete(s.wsClients, fd)
		}
	}
}

func (s *webSocketService) SwitchCitizen(req request.SwitchCitizenReq) error {
	defer func() {
		err := recover()
		if err != nil {
			log.ErrorLogger.Println(err)
		}
	}()

	for _, fd := range s.wsFd[req.CurrentCitizenId] {
		s.wsClients[fd].Conn.WriteJSON(
			response.WebSocketEvent{
				Sender:       "server",
				EventType:    types.WebsocketEventTypeSwitchCitizen,
				EventPayload: req,
			},
		)
	}
	defer func() {
		s.Logout(req.CurrentCitizenId)
		delete(s.wsFd, req.CurrentCitizenId)
	}()

	return nil
}
