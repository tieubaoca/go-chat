package app

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tieubaoca/go-chat-server/db"
	"github.com/tieubaoca/go-chat-server/handlers"
	"github.com/tieubaoca/go-chat-server/middleware"
	"github.com/tieubaoca/go-chat-server/repositories"
	"github.com/tieubaoca/go-chat-server/services"
	"github.com/tieubaoca/go-chat-server/utils/log"
	"go.mongodb.org/mongo-driver/mongo"
)

var database *mongo.Database

func init() {
	log.InfoLogger.Println("Starting app")

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		os.Exit(1)
	}
	log.New(io.MultiWriter(os.Stdout, file))
	gin.DefaultWriter = io.MultiWriter(os.Stdout, file)
	if godotenv.Load() != nil {
		log.FatalLogger.Fatal("Error loading .env file")
	}
}

func Start() {

	database = db.NewDbClient(
		os.Getenv("MONGO_CONNECTION_STRING"),
		os.Getenv("MONGO_DB"),
	)

	saasDb := db.NewSaasMongoDb(
		os.Getenv("SAAS_DB_CONNECTION_STRING"),
		os.Getenv("SAAS_DB"),
	)

	citizenRepository := repositories.NewCitizenRepository(
		saasDb,
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
		citizenRepository,
	)
	messageService := services.NewMessageService(
		messageRepository,
		chatRoomRepository,
	)

	websocketService := services.NewWebSocketService(
		chatRoomRepository,
		messageRepository,
		userRepository,
	)

	saasService := services.NewSaasService()

	authenticationHandler := handlers.NewAuthenticationHandler(websocketService, saasService)
	chatRoomHandler := handlers.NewChatRoomHandler(chatRoomService, saasService)
	messageHandler := handlers.NewMessageHandler(messageService, chatRoomService, saasService)
	userHandler := handlers.NewUserHandler(userService, saasService)
	websocketHandler := handlers.NewWebSocketHandler(websocketService, saasService)

	r := gin.Default()

	limiter := middleware.NewLimiterMiddleware(
		nil,
		10,
		10*time.Second,
	)

	r.Use(limiter.IPRateLimiter())
	r.GET("/saas/api/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

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
		chatRoom.POST("/pagination", chatRoomHandler.PaginationChatRoomBySaId)
		chatRoom.POST("/dm/:member", chatRoomHandler.FindDMByMember)
		// chatRoom.POST("/dm", chatRoomHandler.CreateNewDMChat)
		// chatRoom.POST("/group", chatRoomHandler.CreateNewGroupChat)
		chatRoom.POST("/group/members", chatRoomHandler.FindGroupsByMembers)
		// chatRoom.POST("/group/add-member", chatRoomHandler.AddMemberToGroup)
		// chatRoom.POST("/group/remove-member", chatRoomHandler.RemoveMemberFromGroup)
		// chatRoom.POST("/group/leave/:chatRoomId", chatRoomHandler.LeaveGroup)
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

	switchCitizen := r.Group("/saas/api/switch-citizen")
	switchCitizen.Use(middleware.WhitelistIPsMiddleware())
	switchCitizen.POST("/", websocketHandler.SwitchCitizen)

	// r.GET("/").Handler(http.FileServer(http.Dir("./public")))
	port := os.Getenv("SERVER_PORT")
	log.InfoLogger.Println("Server start on " + port)
	log.FatalLogger.Fatal(http.ListenAndServe(":"+port, r))

	defer database.Client().Disconnect(context.TODO())
}
