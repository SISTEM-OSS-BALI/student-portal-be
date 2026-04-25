package statementletteraiapprovals

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

	protected.POST("/statement-letter-ai-approvals", handler.CreateOrUpdate)
	protected.GET("/statement-letter-ai-approvals", handler.List)
	protected.GET("/statement-letter-ai-approvals/:id", handler.GetByID)
	protected.PUT("/statement-letter-ai-approvals/:id", handler.Update)
}
