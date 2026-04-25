package generatestatementletterai

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

	protected.POST("/generate-statement-letter-ai", handler.Generate)
	protected.POST("/generate-statement-letter-ai/documents", documentHandler.Upsert)
	protected.POST("/generate-statement-letter-ai/documents/:id/submit-to-director", documentHandler.SubmitToDirector)
	protected.POST("/generate-statement-letter-ai/documents/:id/cancel-submit-to-director", documentHandler.CancelSubmitToDirector)
	protected.GET("/generate-statement-letter-ai/documents", documentHandler.List)
	protected.GET("/generate-statement-letter-ai/template", documentHandler.Template)
	protected.GET("/generate-statement-letter-ai/documents/:id/download-pdf", documentHandler.DownloadPDF)
	protected.GET("/generate-statement-letter-ai/documents/student/:student_id", documentHandler.GetByStudentID)
}
