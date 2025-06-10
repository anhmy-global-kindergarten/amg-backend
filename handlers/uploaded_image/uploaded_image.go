package uploaded_image

import (
	"amg-backend/config"
	"amg-backend/models"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"path/filepath"
	"time"
)

// UploadContentImage godoc
// @Summary Upload an image for content
// @Description Uploads an image, saves it, and returns its public URL. Marks the image as 'pending'.
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image file to upload"
// @Success 200 {object} map[string]string{url=string}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/images/upload-image [post]
func (h *UploadedImageHandler) UploadContentImage(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Image file is required"})
	}

	// Tạo tên file duy nhất
	ext := filepath.Ext(file.Filename)
	uniqueFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	savePath := fmt.Sprintf("./uploads/%s", uniqueFilename)

	// Lưu file vật lý
	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save image file"})
	}

	// Tạo URL công khai
	imageURL := fmt.Sprintf("%s/uploads/%s", config.BaseURL, uniqueFilename)

	// Tạo bản ghi trong database
	imageRecord := models.UploadedImage{
		ID:        primitive.NewObjectID(),
		Filename:  uniqueFilename,
		Path:      savePath,
		URL:       imageURL,
		Status:    models.ImageStatusPending,
		CreatedAt: time.Now(),
	}

	collection := h.DB.Database(config.DBName).Collection("UploadedImage")
	_, err = collection.InsertOne(context.TODO(), &imageRecord)
	if err != nil {
		// Nếu không lưu được vào DB, cần cân nhắc xóa file đã lưu
		os.Remove(savePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to record image metadata"})
	}

	return c.JSON(fiber.Map{
		"url": imageURL,
	})
}
