package documents

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

	protected.POST("/documents", handler.Create)
	protected.GET("/documents", handler.List)
	protected.GET("/documents/:id", handler.GetByID)
	protected.PUT("/documents/:id", handler.Update)
	protected.DELETE("/documents/:id", handler.Delete)
	protected.GET("/documents/translations-required", handler.DocumentTranslationRequired)
	protected.POST("/documents/page-count", handler.CountPDFPages)
	protected.POST("/documents/merge-pdf", handler.MergePDF)
}
