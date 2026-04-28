package main

import (
	"log"

	"github.com/username/gin-gorm-api/internal/config"
	"github.com/username/gin-gorm-api/internal/db"
	"github.com/username/gin-gorm-api/internal/modules"
	"github.com/username/gin-gorm-api/internal/router"
	"github.com/username/gin-gorm-api/internal/schema"
)

func main() {
	config.LoadEnv()

	// connect DB
	db.ConnectDB()

	// run versioned migrations
	if err := schema.Migrate(db.DB); err != nil {
		log.Fatalf("failed migrate: %v", err)
	}



	r := router.SetupRouter()
	api := r.Group("/api")
	modules.RegisterAll(api, db.DB)
	r.Run(":8080")
}
