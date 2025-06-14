package uploaded_image

import (
	"amg-backend/config"
	"amg-backend/models"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
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
	imageURL := fmt.Sprintf("/uploads/%s", uniqueFilename)

	// Tạo bản ghi trong database
	imageRecord := models.UploadedImage{
		ID:        primitive.NewObjectID(),
		Filename:  uniqueFilename,
		Path:      savePath,
		URL:       imageURL,
		Status:    models.ImageStatusPending,
		Style:     "",
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

type UpdateImagePayload struct {
	URL   string `json:"url"`
	Style string `json:"style"`
}

// UpdateImagesStatus godoc
// @Summary Update status and style of multiple images
// @Description Receives an array of image URLs and their styles, marks them as 'used' and saves their styles.
// @Tags upload
// @Accept json
// @Produce json
// @Param images body []UpdateImagePayload true "Array of images to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/images/update-status [post]
func (h *UploadedImageHandler) UpdateImagesStatus(c *fiber.Ctx) error {
	var payloads []UpdateImagePayload
	if err := c.BodyParser(&payloads); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if len(payloads) == 0 {
		return c.JSON(fiber.Map{"message": "No images to update"})
	}

	collection := h.DB.Database(config.DBName).Collection("UploadedImage")

	for _, payload := range payloads {
		filter := bson.M{"url": payload.URL}
		update := bson.M{
			"$set": bson.M{
				"status": models.ImageStatusUsed,
				"style":  payload.Style,
			},
		}
		_, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			// Log lỗi nhưng tiếp tục xử lý các ảnh khác
			fmt.Printf("Warning: could not update image %s: %v\n", payload.URL, err)
		}
	}

	return c.JSON(fiber.Map{"message": "Images updated successfully"})
}
