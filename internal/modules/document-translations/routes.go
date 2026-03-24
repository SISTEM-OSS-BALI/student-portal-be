package documenttranslations

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

	protected.POST("/document-translations", handler.Create)
	protected.GET("/document-translations", handler.List)
	protected.GET("/document-translations/:id", handler.GetByID)
	protected.PUT("/document-translations/:id", handler.Update)
	protected.DELETE("/document-translations/:id", handler.Delete)
}
