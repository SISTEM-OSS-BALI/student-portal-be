package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/username/gin-gorm-api/internal/notify"
	"gorm.io/gorm"

	"github.com/username/gin-gorm-api/internal/modules/user"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	repo := user.NewRepository(db)
	notifier := notify.NewService(db)
	service := NewService(repo, notifier, db)
	handler := NewHandler(service)

	rg.POST("/auth/login", handler.Login)
	rg.POST("/auth/register", handler.Register)
	rg.POST("/auth/forgot-password", handler.ForgotPassword)
	rg.POST("/auth/reset-password", handler.ResetPassword)

	protected := rg.Group("")
	protected.Use(AuthMiddleware())
	protected.GET("/auth/me", handler.Me)
}
