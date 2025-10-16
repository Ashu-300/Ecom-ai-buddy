package services

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// Cloudinary client (global)
var Cloud *cloudinary.Cloudinary

func CloudinaryInit() {
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		log.Fatal("CLOUDINARY_URL environment variable is not set")
	}

	var err error
	Cloud, err = cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}

	log.Println("✅ Cloudinary client initialized successfully")
}

// UploadImage uploads an image (from []byte) to Cloudinary
func UploadImage(file io.Reader) (*uploader.UploadResult, error) {
    ctx := context.Background()

    res, err := Cloud.Upload.Upload(ctx, file, uploader.UploadParams{
        Folder: "products",
    })
    if err != nil {
        return nil, err
    }

    return res, nil
}
