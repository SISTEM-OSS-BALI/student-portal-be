package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/username/gin-gorm-api/internal/modules/user"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	repo := user.NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	rg.POST("/auth/login", handler.Login)
	rg.POST("/auth/register", handler.Register)

	protected := rg.Group("")
	protected.Use(AuthMiddleware())
	protected.GET("/auth/me", handler.Me)
}
