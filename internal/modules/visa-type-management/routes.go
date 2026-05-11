package visatype

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

	protected.POST("/visa-types", handler.Create)
	protected.GET("/visa-types", handler.List)
	protected.GET("/visa-types/:id", handler.GetByID)
	protected.PUT("/visa-types/:id", handler.Update)
	protected.DELETE("/visa-types/:id", handler.Delete)
}

