package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	rg.POST("/users", handler.Create)
	rg.GET("/users", handler.List)
	rg.GET("/users/role-students", handler.ListStudents)
	rg.GET("/users/:id", handler.GetByID)
}

func RegisterProtectedRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	rg.PUT("/users/:id", handler.Update)
	rg.PATCH("/users/:id/translation-quota", handler.PatchQuotaTranslation)
	rg.PATCH("/users/:id/visa-status", handler.PatchVisaStatus)
	rg.PATCH("/users/:id/student-status", handler.PatchStudentStatus)
	rg.DELETE("/users/:id", handler.Delete)
	rg.PATCH("/users/:id/document-consent", handler.PatchDocumentConsent)
}
