package app

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tieubaoca/go-chat-server/controllers"
	"github.com/tieubaoca/go-chat-server/middleware"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/utils/log"
)

func Start() {

	r := gin.Default()
	authorized := r.Group("/saas/api")
	authorized.Use(middleware.JwtMiddleware)
	authorized.GET("/ws", controllers.HandleWebSocket)

	authorized.POST("/auth", controllers.Authentication)
	r.POST("/saas/api/get-access-token", controllers.GetAccessToken)

	authorized.GET("/chat-room/id/{id}", controllers.FindChatRoomById)
	authorized.GET("/chat-room/all", controllers.FindChatRooms)
	authorized.POST("/chat-room/dm/members", controllers.FindDMByMembers)
	authorized.POST("/chat-room/dm", controllers.CreateDMRoom)
	authorized.POST("/chat-room/group", controllers.CreateNewGroupChat)
	authorized.POST("/chat-room/group/members", controllers.FindGroupsByMembers)

	authorized.POST("/user/online/pagination", controllers.PaginationOnlineFriend)
	authorized.POST("/saas/api/user/logout", controllers.Logout)

	authorized.GET("/message/chat-room/{chatRoomId}", controllers.FindMessagesByChatRoomId)
	authorized.POST("/message/pagination", controllers.PaginationMessagesByChatRoomId)

	// r.GET("/").Handler(http.FileServer(http.Dir("./public")))
	port := os.Getenv("SERVER_PORT")
	log.InfoLogger.Println("Server start on " + port)
	log.FatalLogger.Fatal(http.ListenAndServe(":"+port, r))

	defer services.GetDBClient().Client().Disconnect(context.TODO())
}
