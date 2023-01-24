package app

import (
	"context"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/tieubaoca/go-chat-server/controllers"
	"github.com/tieubaoca/go-chat-server/middleware"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/types"
	"github.com/tieubaoca/go-chat-server/utils/log"
)

func Start() {

	r := mux.NewRouter().StrictSlash(true)
	r.Handle("/saas/api/ws", http.HandlerFunc(middleware.JwtMiddleware(controllers.HandleWebSocket, types.RoleViewProfile)))

	r.HandleFunc("/saas/api/auth", middleware.JwtMiddleware(controllers.Authentication, types.RoleViewProfile)).Methods("POST")
	r.HandleFunc("/saas/api/get-access-token", controllers.GetAccessToken).Methods("POST")

	r.HandleFunc("/saas/api/chat-room/id/{id}", middleware.JwtMiddleware(controllers.FindChatRoomById, types.RoleViewProfile)).Methods("GET")
	r.HandleFunc("/saas/api/chat-room/all", middleware.JwtMiddleware(controllers.FindChatRooms, types.RoleViewProfile)).Methods("GET")
	r.HandleFunc("/saas/api/chat-room/group", middleware.JwtMiddleware(controllers.CreateNewGroupChat, types.RoleManagerAccount)).Methods("POST")
	r.HandleFunc("/saas/api/chat-room/dm", middleware.JwtMiddleware(controllers.CreateDMRoom, types.RoleManagerAccount)).Methods("POST")
	r.HandleFunc("/saas/api/chat-room/dm/members", middleware.JwtMiddleware(controllers.FindDMByMembers, types.RoleViewProfile)).Methods("POST")
	r.HandleFunc("/saas/api/chat-room/group/members", middleware.JwtMiddleware(controllers.FindGroupsByMembers, types.RoleViewProfile)).Methods("POST")

	r.HandleFunc("/saas/api/user/online/pagination", middleware.JwtMiddleware(controllers.PaginationOnlineFriend, types.RoleViewProfile)).Methods("POST")
	r.HandleFunc("/saas/api/user/logout", middleware.JwtMiddleware(controllers.Logout, types.RoleManagerAccount)).Methods("POST")

	r.HandleFunc("/saas/api/message/chat-room/{chatRoomId}", middleware.JwtMiddleware(controllers.FindMessagesByChatRoomId, types.RoleViewProfile)).Methods("GET")
	r.HandleFunc("/saas/api/message/pagination", middleware.JwtMiddleware(controllers.PaginationMessagesByChatRoomId, types.RoleViewProfile)).Methods("POST")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))
	port := os.Getenv("SERVER_PORT")
	log.InfoLogger.Println("Server start on " + port)
	log.FatalLogger.Fatal(http.ListenAndServe(":"+port, r))

	defer services.GetDBClient().Client().Disconnect(context.TODO())
}
