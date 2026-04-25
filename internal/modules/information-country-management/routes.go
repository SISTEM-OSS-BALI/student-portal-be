package informationcountrymanagement

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

	protected.POST("/information-countries", handler.Create)
	protected.GET("/information-countries", handler.List)
	protected.GET("/information-countries/:id", handler.GetByID)
	protected.GET("/information-countries/slug/:slug", handler.GetBySlug)
	protected.GET("/information-countries/country/:country_id", handler.ListByCountryID)
	protected.PUT("/information-countries/:id", handler.Update)
	protected.DELETE("/information-countries/:id", handler.Delete)
}
