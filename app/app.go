package app

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tieubaoca/go-chat-server/controllers"
	"github.com/tieubaoca/go-chat-server/middleware"
	"github.com/tieubaoca/go-chat-server/services"
)

var HomeDir string

func Start() {
	r := mux.NewRouter().StrictSlash(true)
	r.Handle("/api/ws", http.HandlerFunc(middleware.JwtMiddleware(controllers.HandleWebSocket)))

	r.HandleFunc("/api/auth", middleware.JwtMiddleware(controllers.Authentication)).Methods("POST")
	r.HandleFunc("/api/get-access-token", controllers.GetAccessToken).Methods("POST")

	r.HandleFunc("/api/chat-room/id/{id}", controllers.FindChatRoomById).Methods("GET")
	r.HandleFunc("/api/chat-room/member/{member}", controllers.FindChatRoomsByMember).Methods("GET")
	r.HandleFunc("/api/chat-room/members", controllers.FindDMByMembers).Methods("POST")
	r.HandleFunc("/api/chat-room", controllers.InsertChatRoom).Methods("POST")

	r.HandleFunc("/api/user/username/{username}", controllers.FindUserByUsername).Methods("GET")
	r.HandleFunc("/api/user/id/{id}", controllers.FindUserById).Methods("GET")
	r.HandleFunc("/api/user/online", controllers.FindOnlineUsers).Methods("GET")

	r.HandleFunc("/api/message/chat-room/{chatRoomId}", controllers.FindMessagesByChatRoomId).Methods("GET")
	r.HandleFunc("/api/message/pagination", controllers.PaginationMessagesByChatRoomId).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))

	log.Println("Server start on 8800")

	log.Fatal(http.ListenAndServe(":8800", r))

	defer services.GetDBClient().Client().Disconnect(context.TODO())
}
