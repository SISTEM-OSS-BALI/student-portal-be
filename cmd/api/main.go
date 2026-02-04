package main

import (
	"log"

	"github.com/username/gin-gorm-api/internal/config"
	"github.com/username/gin-gorm-api/internal/db"
	"github.com/username/gin-gorm-api/internal/handler"
	"github.com/username/gin-gorm-api/internal/models"
	"github.com/username/gin-gorm-api/internal/repository"
	"github.com/username/gin-gorm-api/internal/router"
	"github.com/username/gin-gorm-api/internal/service"
)

func main() {
	config.LoadEnv()

	// connect DB
	db.ConnectDB()

	// migrate
	if err := db.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("failed migrate: %v", err)
	}

	userRepo := repository.NewGormUserRepository(db.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r := router.SetupRouter(userHandler)
	r.Run(":8080")
}
