package app

import (
	"context"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/tieubaoca/go-chat-server/controllers"
	"github.com/tieubaoca/go-chat-server/middleware"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/types"
)

func Start() {

	r := mux.NewRouter().StrictSlash(true)
	r.Handle("/saas/api/ws", http.HandlerFunc(middleware.JwtMiddleware(controllers.HandleWebSocket, types.RoleViewProfile)))

	r.HandleFunc("/saas/api/auth", middleware.JwtMiddleware(controllers.Authentication, types.RoleViewProfile)).Methods("POST")
	r.HandleFunc("/saas/api/get-access-token", controllers.GetAccessToken).Methods("POST")

	// r.HandleFunc("/saas/api/chat-room/id/{id}", controllers.FindChatRoomById).Methods("GET")
	r.HandleFunc("/saas/api/chat-room/member/{member}", middleware.JwtMiddleware(controllers.FindChatRooms, types.RoleViewProfile)).Methods("GET")
	r.HandleFunc("/saas/api/chat-room/group", middleware.JwtMiddleware(controllers.CreateNewGroupChat, types.RoleManagerAccount)).Methods("POST")
	r.HandleFunc("/saas/api/chat-room/dm", middleware.JwtMiddleware(controllers.CreateDMRoom, types.RoleManagerAccount)).Methods("POST")
	r.HandleFunc("/saas/api/chat-room/dm/members", middleware.JwtMiddleware(controllers.FindDMByMembers, types.RoleViewProfile)).Methods("POST")
	r.HandleFunc("/saas/api/chat-room/group/members", middleware.JwtMiddleware(controllers.FindGroupsByMembers, types.RoleViewProfile)).Methods("POST")

	// r.HandleFunc("/saas/api/user/username/{username}", controllers.FindUserByUsername).Methods("GET")
	// r.HandleFunc("/saas/api/user/id/{id}", controllers.FindUserById).Methods("GET")
	r.HandleFunc("/saas/api/user/online", middleware.JwtMiddleware(controllers.FindOnlineFriends, types.RoleViewProfile)).Methods("POST")
	r.HandleFunc("/saas/api/user/logout", middleware.JwtMiddleware(controllers.Logout, types.RoleManagerAccount)).Methods("POST")

	r.HandleFunc("/saas/api/message/chat-room/{chatRoomId}", middleware.JwtMiddleware(controllers.FindMessagesByChatRoomId, types.RoleViewProfile)).Methods("GET")
	r.HandleFunc("/saas/api/message/pagination", middleware.JwtMiddleware(controllers.PaginationMessagesByChatRoomId, types.RoleViewProfile)).Methods("POST")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))

	port := os.Getenv("SERVER_PORT")
	log.Info("Server start on " + port)

	log.Fatal(http.ListenAndServe(":"+port, r))

	defer services.GetDBClient().Client().Disconnect(context.TODO())
}
