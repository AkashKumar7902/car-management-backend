package routes

import (
	"github.com/akashkumar7902/car-management-backend/config"
	"github.com/akashkumar7902/car-management-backend/controllers"
	"github.com/akashkumar7902/car-management-backend/middlewares"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CarRoutes(r *gin.Engine, db *gorm.DB, cfg config.Config, cld *cloudinary.Cloudinary) {
	carController := controllers.CarController{
		DB:         db,
		Cfg:        cfg,
		Cloudinary: cld,
	}

	// Apply authentication middleware
	authMiddleware := middlewares.AuthMiddleware(db, cfg)

	cars := r.Group("/api/cars").Use(authMiddleware)
	{
		cars.POST("/", carController.CreateCar)
		cars.GET("/", carController.ListCars)
		cars.GET("/search", carController.SearchCars)
		cars.GET("/:id", carController.GetCar)
		cars.PUT("/:id", carController.UpdateCar)
		cars.DELETE("/:id", carController.DeleteCar)
	}
}
