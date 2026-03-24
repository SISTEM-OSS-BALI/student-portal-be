package country

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

	protected.POST("/countries", handler.Create)
	protected.GET("/countries", handler.List)
	protected.GET("/countries/:id", handler.GetByID)
	protected.PUT("/countries/:id", handler.Update)
	protected.DELETE("/countries/:id", handler.Delete)
}
