package uploaded_image

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UploadedImageHandler struct {
	Router fiber.Router
	DB     *mongo.Client
}

func RegisterUploadedImageHandler(router fiber.Router, db *mongo.Client) {
	uploadedImageHandler := UploadedImageHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Post("/upload-image", uploadedImageHandler.UploadContentImage)
}
