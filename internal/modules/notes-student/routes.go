package notesstudent

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

	protected.POST("/notes-student", handler.Create)
	protected.GET("/notes-student", handler.List)
	protected.GET("/notes-student/:id", handler.GetByID)
	protected.PUT("/notes-student/:id", handler.Update)
	protected.DELETE("/notes-student/:id", handler.Delete)
}
