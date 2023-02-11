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
	r.POST("/saas/api/login", controllers.GetAccessToken)

	authorized := r.Group("/saas/api")
	authorized.Use(middleware.JwtMiddleware)
	authorized.GET("/ws", controllers.HandleWebSocket)
	authorized.POST("/auth", controllers.Authentication)

	chatRoom := r.Group("/saas/api/chat-room")
	chatRoom.Use(middleware.JwtMiddleware)
	{
		chatRoom.GET("/id/:id", controllers.FindChatRoomById)
		chatRoom.GET("/all", controllers.FindChatRooms)
		chatRoom.POST("/dm/members", controllers.FindDMByMembers)
		chatRoom.POST("/dm", controllers.CreateDMRoom)
		chatRoom.POST("/group", controllers.CreateNewGroupChat)
		chatRoom.POST("/group/members", controllers.FindGroupsByMembers)
		chatRoom.POST("/group/add-member", controllers.AddMemberToGroup)
		chatRoom.POST("/group/remove-member", controllers.RemoveMemberFromGroup)
		chatRoom.POST("/group/leave/:chatRoomId", controllers.LeaveGroup)
	}

	user := r.Group("/saas/api/user")
	user.Use(middleware.JwtMiddleware)
	{
		user.POST("/online/pagination", controllers.PaginationOnlineFriend)
		user.POST("/logout", controllers.Logout)
	}

	message := r.Group("/saas/api/message")
	message.Use(middleware.JwtMiddleware)
	{
		message.GET("/chat-room/:chatRoomId", controllers.FindMessagesByChatRoomId)
		message.POST("/pagination", controllers.PaginationMessagesByChatRoomId)
	}

	// r.GET("/").Handler(http.FileServer(http.Dir("./public")))
	port := os.Getenv("SERVER_PORT")
	log.InfoLogger.Println("Server start on " + port)
	log.FatalLogger.Fatal(http.ListenAndServe(":"+port, r))

	defer services.GetDBClient().Client().Disconnect(context.TODO())
}
