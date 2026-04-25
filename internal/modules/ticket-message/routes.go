package ticketmessage

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/username/gin-gorm-api/internal/modules/auth"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	protected := rg.Group("")
	protected.Use(auth.AuthMiddleware())

	protected.POST("/ticket-messages", handler.Create)
	protected.GET("/ticket-messages", handler.List)
	protected.GET("/ticket-messages/:id", handler.GetByID)
	protected.PUT("/ticket-messages/:id", handler.Update)
	protected.PUT("/ticket-messages/:id/status", handler.UpdateStatus)
	protected.DELETE("/ticket-messages/:id/with-conversation", handler.DeleteWithConversation)
	protected.DELETE("/ticket-messages/:id", handler.Delete)
}
