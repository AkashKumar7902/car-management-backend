package config

import (
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
)

func InitCloudinary(cfg Config) *cloudinary.Cloudinary {
	cld, err := cloudinary.NewFromParams(cfg.CloudName, cfg.CloudAPIKey, cfg.CloudAPISecret)
	if err != nil {
		log.Fatal("Failed to initialize Cloudinary:", err)
	}
	return cld
}
