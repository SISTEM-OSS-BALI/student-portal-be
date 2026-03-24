package answerapprovals

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

	protected.POST("/answer-approvals", handler.CreateOrUpdate)
	protected.GET("/answer-approvals", handler.List)
	protected.GET("/answer-approvals/:id", handler.GetByID)
	protected.PUT("/answer-approvals/:id", handler.Update)
	protected.DELETE("/answer-approvals/:id", handler.Delete)
}
