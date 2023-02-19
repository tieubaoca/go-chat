package app

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tieubaoca/go-chat-server/db"
	"github.com/tieubaoca/go-chat-server/handlers"
	"github.com/tieubaoca/go-chat-server/middleware"
	"github.com/tieubaoca/go-chat-server/repositories"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/utils/log"
	"go.mongodb.org/mongo-driver/mongo"
)

var database *mongo.Database

func Start() {

	database = db.NewDbClient(
		os.Getenv("MONGO_CONNECTION_STRING"),
		os.Getenv("MONGO_DB"),
	)

	userRepository := repositories.NewUserRepository(
		database,
	)
	chatRoomRepository := repositories.NewChatRoomRepository(
		database,
	)
	messageRepository := repositories.NewMessageRepository(
		database,
	)

	userService := services.NewUserService(
		userRepository,
	)
	chatRoomService := services.NewChatRoomService(
		chatRoomRepository,
		userRepository,
	)
	messageService := services.NewMessageService(
		messageRepository,
	)

	websocketService := services.NewWebSocketService(
		chatRoomRepository,
		messageRepository,
		userRepository,
	)

	authenticationHandler := handlers.NewAuthenticationHandler(websocketService)
	chatRoomHandler := handlers.NewChatRoomHandler(chatRoomService)
	messageHandler := handlers.NewMessageHandler(messageService, chatRoomService)
	userHandler := handlers.NewUserHandler(userService)
	websocketHandler := handlers.NewWebSocketHandler(websocketService)

	r := gin.Default()
	r.POST("/saas/api/login", authenticationHandler.Login)

	authentication := r.Group("/saas/api")
	{
		authentication.Use(middleware.JwtMiddleware)
		authentication.GET("/ws", websocketHandler.HandleWebSocket)
		authentication.POST("/logout", authenticationHandler.Logout)
	}

	chatRoom := r.Group("/saas/api/chat-room")
	chatRoom.Use(middleware.JwtMiddleware)
	{
		chatRoom.GET("/id/:id", chatRoomHandler.FindChatRoomById)
		chatRoom.GET("/all", chatRoomHandler.FindChatRoomsBySaId)
		chatRoom.POST("/dm/:member", chatRoomHandler.FindDMByMember)
		// chatRoom.POST("/dm", chatRoomHandler.CreateNewDMChat)
		// chatRoom.POST("/group", chatRoomHandler.CreateNewGroupChat)
		chatRoom.POST("/group/members", chatRoomHandler.FindGroupsByMembers)
		chatRoom.POST("/group/add-member", chatRoomHandler.AddMemberToGroup)
		chatRoom.POST("/group/remove-member", chatRoomHandler.RemoveMemberFromGroup)
		chatRoom.POST("/group/leave/:chatRoomId", chatRoomHandler.LeaveGroup)
	}

	user := r.Group("/saas/api/user")
	user.Use(middleware.JwtMiddleware)
	{
		user.POST("/online/pagination", userHandler.PaginationOnlineFriend)
	}

	message := r.Group("/saas/api/message")
	message.Use(middleware.JwtMiddleware)
	{
		message.GET("/chat-room/:chatRoomId", messageHandler.FindMessagesByChatRoomId)
		message.POST("/pagination", messageHandler.PaginationMessagesByChatRoomId)
	}

	// r.GET("/").Handler(http.FileServer(http.Dir("./public")))
	port := os.Getenv("SERVER_PORT")
	log.InfoLogger.Println("Server start on " + port)
	log.FatalLogger.Fatal(http.ListenAndServe(":"+port, r))

	defer database.Client().Disconnect(context.TODO())
}
