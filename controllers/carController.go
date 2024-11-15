package controllers

import (
	"log"
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

// CreateCar handles creating a new car with optional image uploads
// @Summary Create a new car
// @Description Create a new car with title, description, tags, and optional images
// @Tags Cars
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Title"
// @Param description formData string false "Description"
// @Param tags formData string false "Tags (comma-separated)"
// @Param images formData file false "Images" maxItems(10)
// @Success 201 {object} models.Car
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /api/cars [post]
func (cc *CarController) CreateCar(c *gin.Context) {
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

	car := models.Car{
		UserID:      user.ID,
		Title:       title,
		Description: description,
		Tags:        tagList,
		Images:      []string{}, // Initialize as empty slice
	}

	// Handle image uploads
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	files := form.File["images"]
	for _, file := range files {
		// Upload to Cloudinary
		uploadParams := uploader.UploadParams{
			Folder: "car_management",
		}

		uploadResult, err := cc.Cloudinary.Upload.Upload(c, file, uploadParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Image upload failed"})
			return
		}
		car.Images = append(car.Images, uploadResult.SecureURL)
	}

	if err := cc.DB.Create(&car).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create car"})
		return
	}

	c.JSON(http.StatusCreated, car)
}

// ListCars lists all cars of the logged-in user
// @Summary List all cars
// @Description Get a list of all cars for the logged-in user
// @Tags Cars
// @Accept json
// @Produce json
// @Success 200 {array} models.Car
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /api/cars [get]
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
// @Summary Get a specific car
// @Description Get car by ID for the logged-in user
// @Tags Cars
// @Accept json
// @Produce json
// @Param id path int true "Car ID"
// @Success 200 {object} models.Car
// @Failure 401 {object} error
// @Failure 403 {object} error
// @Failure 404 {object} error
// @Router /api/cars/{id} [get]
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
// @Summary Update a car
// @Description Update car details for the logged-in user
// @Tags Cars
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Car ID"
// @Param title formData string false "Title"
// @Param description formData string false "Description"
// @Param tags formData string false "Tags (comma-separated)"
// @Param images formData string false "Images" maxItems(10)
// @Success 200 {object} models.Car
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 403 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /api/cars/{id} [put]
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

	// Parse form data
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	// Get form fields
	title := c.PostForm("title")
	description := c.PostForm("description")
	tagsStr := c.PostForm("tags")

	// Update car fields if provided
	if title != "" {
		car.Title = title
	}
	if description != "" {
		car.Description = description
	}
	if tagsStr != "" {
		tags := strings.Split(tagsStr, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
		car.Tags = tags
	}

	// Handle image uploads
	form := c.Request.MultipartForm
	files := form.File["images"]
	if len(files) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 10 images allowed"})
		return
	}

	newImageUrls := []string{}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open image"})
			return
		}
		defer file.Close()

		// Upload to Cloudinary
		uploadResult, err := cc.Cloudinary.Upload.Upload(c, file, uploader.UploadParams{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
			return
		}

		newImageUrls = append(newImageUrls, uploadResult.SecureURL)
	}

	// Append new images to existing images
	if len(newImageUrls) > 0 {
		car.Images = append(car.Images, newImageUrls...)
	}

	// Save updated car
	if err := cc.DB.Save(&car).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update car"})
		return
	}

	c.JSON(http.StatusOK, car)
}

// DeleteCar deletes a specific car
// @Summary Delete a car
// @Description Delete a car by ID for the logged-in user
// @Tags Cars
// @Accept json
// @Produce json
// @Param id path int true "Car ID"
// @Success 200 {object} error
// @Failure 401 {object} error
// @Failure 403 {object} error
// @Failure 404 {object} error
// @Router /api/cars/{id} [delete]
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
// @Summary Search cars
// @Description Search cars by keyword in title, description, or tags
// @Tags Cars
// @Accept json
// @Produce json
// @Param keyword query string true "Search keyword"
// @Success 200 {array} models.Car
// @Failure 400 {object} error
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Router /api/cars/search [get]
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
    if err := cc.DB.Where(
        "user_id = ? AND (title ILIKE ? OR description ILIKE ? OR ? = ANY(tags))",
        user.ID, "%"+keyword+"%", "%"+keyword+"%", keyword,
    ).Find(&cars).Error; err != nil {
        log.Println(err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search cars"})
        return
    }

    c.JSON(http.StatusOK, cars)
}
