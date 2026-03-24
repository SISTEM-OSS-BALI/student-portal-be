package steps

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

	protected.POST("/steps", handler.Create)
	protected.GET("/steps", handler.List)
	protected.GET("/steps/:id", handler.GetByID)
	protected.PUT("/steps/:id", handler.Update)
	protected.DELETE("/steps/:id", handler.Delete)
}
