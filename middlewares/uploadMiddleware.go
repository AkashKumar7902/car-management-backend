package middlewares

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func UploadMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Maximum number of images
		const MaxImages = 10

		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart form"})
			c.Abort()
			return
		}

		files := form.File["images"]
		if len(files) > MaxImages {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 10 images allowed"})
			c.Abort()
			return
		}

		var imagePaths []string
		for _, file := range files {
			// Validate file type
			if !isImage(file) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPEG, JPG, and PNG images are allowed"})
				c.Abort()
				return
			}

			// Save file
			filename := fmt.Sprintf("%d_%s", c.Writer.Status(), filepath.Base(file.Filename))
			filepath := filepath.Join("uploads", filename)

			if err := c.SaveUploadedFile(file, filepath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
				c.Abort()
				return
			}

			imagePaths = append(imagePaths, filepath)
		}

		// Set image paths in context
		c.Set("imagePaths", imagePaths)

		c.Next()
	}
}

func isImage(file *multipart.FileHeader) bool {
	allowedTypes := []string{".jpeg", ".jpg", ".png"}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	for _, t := range allowedTypes {
		if t == ext {
			return true
		}
	}
	return false
}
