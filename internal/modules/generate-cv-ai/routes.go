package generatecvai

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/username/gin-gorm-api/internal/modules/auth"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	service := NewService()
	handler := NewHandler(service)
	documentRepo := NewGeneratedDocumentRepository(db)
	documentService := NewGeneratedDocumentService(documentRepo)
	documentHandler := NewGeneratedDocumentHandler(documentService)

	protected := rg.Group("")
	protected.Use(auth.AuthMiddleware())

	protected.POST("/generate-cv-ai", handler.Generate)
	protected.POST("/generate-cv-ai/documents", documentHandler.Upsert)
	protected.GET("/generate-cv-ai/documents", documentHandler.List)
	protected.GET("/generate-cv-ai/documents/:student_id", documentHandler.GetByStudentID)
}
