package controllers

import (
	"net/http"
	"strings"

	"github.com/akashkumar7902/car-management-backend/config"
	"github.com/akashkumar7902/car-management-backend/models"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CarController struct {
	DB         *gorm.DB
	Cfg        config.Config
	Cloudinary *cloudinary.Cloudinary
}

// CreateCar handles creating a new car
func (cc *CarController) CreateCar(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := userInterface.(models.User)

	var input struct {
		Title       string   `form:"title" binding:"required"`
		Description string   `form:"description"`
		Tags        string   `form:"tags"`
		Images      []string `form:"images"` // URLs or paths
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tags := []string{}
	if input.Tags != "" {
		tags = strings.Split(input.Tags, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	}

	car := models.Car{
		UserID:      user.ID,
		Title:       input.Title,
		Description: input.Description,
		Tags:        tags,
		Images:      input.Images,
	}

	if err := cc.DB.Create(&car).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create car"})
		return
	}

	c.JSON(http.StatusCreated, car)
}

// ListCars lists all cars of the logged-in user
func (cc *CarController) ListCars(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := userInterface.(models.User)

	var cars []models.Car
	if err := cc.DB.Where("user_id = ?", user.ID).Find(&cars).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cars"})
		return
	}

	c.JSON(http.StatusOK, cars)
}

// GetCar retrieves a specific car
func (cc *CarController) GetCar(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := userInterface.(models.User)

	carID := c.Param("id")
	var car models.Car
	if err := cc.DB.First(&car, carID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
		return
	}

	if car.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, car)
}

// UpdateCar updates a specific car
func (cc *CarController) UpdateCar(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := userInterface.(models.User)

	carID := c.Param("id")
	var car models.Car
	if err := cc.DB.First(&car, carID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
		return
	}

	if car.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var input struct {
		Title       string   `form:"title"`
		Description string   `form:"description"`
		Tags        string   `form:"tags"`
		Images      []string `form:"images"` // URLs or paths
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Title != "" {
		car.Title = input.Title
	}
	if input.Description != "" {
		car.Description = input.Description
	}
	if input.Tags != "" {
		tags := strings.Split(input.Tags, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
		car.Tags = tags
	}
	if len(input.Images) > 0 {
		car.Images = input.Images
	}

	if err := cc.DB.Save(&car).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update car"})
		return
	}

	c.JSON(http.StatusOK, car)
}

// DeleteCar deletes a specific car
func (cc *CarController) DeleteCar(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := userInterface.(models.User)

	carID := c.Param("id")
	var car models.Car
	if err := cc.DB.First(&car, carID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
		return
	}

	if car.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := cc.DB.Delete(&car).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete car"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Car deleted successfully"})
}

// SearchCars searches cars based on a keyword
func (cc *CarController) SearchCars(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := userInterface.(models.User)

	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Keyword query parameter is required"})
		return
	}

	var cars []models.Car
	searchQuery := "%" + keyword + "%"
	if err := cc.DB.Where("user_id = ? AND (title ILIKE ? OR description ILIKE ? OR tags && ?)", user.ID, searchQuery, searchQuery, []string{keyword}).Find(&cars).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search cars"})
		return
	}

	c.JSON(http.StatusOK, cars)
}

// CreateCarWithCloudinary handles creating a new car with image uploads to Cloudinary
func (cc *CarController) CreateCarWithCloudinary(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := userInterface.(models.User)

	title := c.PostForm("title")
	description := c.PostForm("description")
	tags := c.PostForm("tags")

	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	// Handle tags
	tagList := []string{}
	if tags != "" {
		tagList = strings.Split(tags, ",")
		for i := range tagList {
			tagList[i] = strings.TrimSpace(tagList[i])
		}
	}

	// Handle image uploads
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart form"})
		return
	}

	files := form.File["images"]
	if len(files) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 10 images allowed"})
		return
	}

	imageURLs := []string{}
	for _, file := range files {
		// Open the file
		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image"})
			return
		}
		defer f.Close()

		// Upload to Cloudinary
		uploadParams := uploader.UploadParams{
			Folder: "car_management",
		}

		uploadResult, err := cc.Cloudinary.Upload.Upload(c, f, uploadParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
			return
		}

		imageURLs = append(imageURLs, uploadResult.SecureURL)
	}

	// Create car
	car := models.Car{
		UserID:      user.ID,
		Title:       title,
		Description: description,
		Tags:        tagList,
		Images:      imageURLs,
	}

	if err := cc.DB.Create(&car).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create car"})
		return
	}

	c.JSON(http.StatusCreated, car)
}
