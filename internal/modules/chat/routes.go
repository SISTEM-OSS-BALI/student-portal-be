package chat

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/username/gin-gorm-api/internal/notify"
	"gorm.io/gorm"

	"github.com/username/gin-gorm-api/internal/modules/auth"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	repo := NewRepository(db)
	notifier := notify.NewService(db)
	service := NewService(repo, notifier)

	socketServer, err := NewSocketServer(service)
	if err != nil {
		log.Printf("chat socket server init failed: %v", err)
		return
	}
	handler := NewHandler(service, socketServer)

	protected := rg.Group("")
	protected.Use(auth.AuthMiddleware())

	protected.POST("/chats/conversations", handler.CreateConversation)
	protected.GET("/chats/conversations", handler.ListConversations)
	protected.GET("/chats/conversations/:id/messages", handler.ListMessages)
	protected.POST("/chats/conversations/:id/messages", handler.SendMessage)
	protected.POST("/chats/conversations/:id/read", handler.MarkRead)
	protected.GET("/chats/mentions", handler.ListMentions)
	protected.POST("/chats/mentions/:id/read", handler.MarkMentionRead)

	socketServer.Start()
	rg.Any("/socket.io/*any", gin.WrapH(socketServer.Handler()))
}
