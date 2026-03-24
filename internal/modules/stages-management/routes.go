package stages

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

	protected.POST("/stages", handler.Create)
	protected.GET("/stages", handler.List)
	protected.GET("/stages/:id", handler.GetByID)
	protected.PUT("/stages/:id", handler.Update)
	protected.DELETE("/stages/:id", handler.Delete)
}