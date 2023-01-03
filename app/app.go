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
	r.Handle("/saas/api/ws", http.HandlerFunc(middleware.JwtMiddleware(controllers.HandleWebSocket)))

	r.HandleFunc("/saas/api/auth", middleware.JwtMiddleware(controllers.Authentication)).Methods("POST")
	r.HandleFunc("/saas/api/get-access-token", controllers.GetAccessToken).Methods("POST")

	r.HandleFunc("/saas/api/chat-room/id/{id}", controllers.FindChatRoomById).Methods("GET")
	r.HandleFunc("/saas/api/chat-room/member/{member}", controllers.FindChatRoomsByMember).Methods("GET")
	r.HandleFunc("/saas/api/chat-room/group", controllers.CreateNewGroupChat).Methods("POST")
	r.HandleFunc("/saas/api/chat-room/dm", controllers.CreateDMRoom).Methods("POST")
	r.HandleFunc("/saas/api/chat-room/dm/members", controllers.FindDMByMembers).Methods("POST")
	r.HandleFunc("/saas/api/chat-room/group/members", controllers.FindGroupsByMembers).Methods("POST")

	// r.HandleFunc("/saas/api/user/username/{username}", controllers.FindUserByUsername).Methods("GET")
	// r.HandleFunc("/saas/api/user/id/{id}", controllers.FindUserById).Methods("GET")
	r.HandleFunc("/saas/api/user/online", controllers.FindOnlineUsers).Methods("GET")
	r.HandleFunc("/saas/api/user/logout", controllers.Logout).Methods("POST")

	r.HandleFunc("/saas/api/message/chat-room/{chatRoomId}", controllers.FindMessagesByChatRoomId).Methods("GET")
	r.HandleFunc("/saas/api/message/pagination", controllers.PaginationMessagesByChatRoomId).Methods("POST")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))

	log.Println("Server start on 8800")

	log.Fatal(http.ListenAndServe(":8800", r))

	defer services.GetDBClient().Client().Disconnect(context.TODO())
}
