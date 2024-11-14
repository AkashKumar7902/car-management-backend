package main

import (
	"log"
	"net/http"

	"github.com/akashkumar7902/car-management-backend/config"
	_ "github.com/akashkumar7902/car-management-backend/docs" // Import generated docs
	"github.com/akashkumar7902/car-management-backend/models"
	"github.com/akashkumar7902/car-management-backend/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Car Management API
// @version 1.0
// @description API Documentation for Car Management Application
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

func main() {
	cfg := config.LoadConfig()

	cld := config.InitCloudinary(cfg)

	// Initialize Gin
	r := gin.Default()

	// Serve static files (uploads)
	r.Static("/uploads", "./uploads")

	// Initialize Database
	dsn := "host=" + cfg.DBHost + " user=" + cfg.DBUser + " password=" + cfg.DBPassword + " dbname=" + cfg.DBName + " port=" + cfg.DBPort + " sslmode=disable TimeZone=Asia/Kolkata"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate models
	err = db.AutoMigrate(&models.User{}, &models.Car{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Routes
	routes.AuthRoutes(r, db, cfg)
    routes.CarRoutes(r, db, cfg, cld)

	// Swagger Documentation
	r.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start Server
	if err := r.Run(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to run server: %v", err)
	}
}
