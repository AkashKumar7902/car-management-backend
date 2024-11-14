package routes

import (
	"github.com/akashkumar7902/car-management-backend/config"
	"github.com/akashkumar7902/car-management-backend/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoutes(r *gin.Engine, db *gorm.DB, cfg config.Config) {
	authController := controllers.AuthController{
		DB:  db,
		Cfg: cfg,
	}

	auth := r.Group("/api/users")
	{
		auth.POST("/signup", authController.RegisterUser)
		auth.POST("/login", authController.LoginUser)
	}
}
