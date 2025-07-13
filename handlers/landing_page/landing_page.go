package landing_page

import (
	"amg-backend/config"
	"amg-backend/models"
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// GetLandingPageContent godoc
// @Summary Get Landing Page Content
// @Description Get the content of the landing page
// @Tags landing page
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string "message: Nội dung chưa được khởi tạo"
// @Failure 500 {object} map[string]string "error: Lỗi máy chủ"
// @Router /amg/v1/landing-page/get-content [get]
func (h *LandingPageHandler) GetLandingPageContent(c *fiber.Ctx) error {
	collection := h.DB.Database(config.DBName).Collection("LandingPageContent")

	var result models.LandingPageContent
	filter := bson.M{"key": "main"}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Nội dung chưa được khởi tạo"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}

	return c.JSON(result.Content)
}

// UpdateLandingPageContent godoc
// @Summary Update Landing Page Content
// @Description Update the content of the landing page
// @Tags landing page
// @Accept json
// @Produce json
// @Param content body map[string]interface{} true "Object containing the new content for the landing page"
// @Success 200 {object} map[string]string "message: "Nội dung landing page đã được cập nhật thành công"
// @Failure 400 {object} map[string]string "error: Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "error: Lỗi máy chủ"
// @Router /amg/v1/landing-page/update-content [post]
func (h *LandingPageHandler) UpdateLandingPageContent(c *fiber.Ctx) error {
	var newContent map[string]interface{}

	if err := c.BodyParser(&newContent); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON data"})
	}

	if len(newContent) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No update data provided"})
	}

	collection := h.DB.Database(config.DBName).Collection("LandingPageContent")

	filter := bson.M{"key": "main"}
	update := bson.M{
		"$set": bson.M{
			"content":    newContent,
			"updated_at": time.Now(),
		},
		"$setOnInsert": bson.M{
			"key": "main",
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Update failed in database"})
	}

	return c.JSON(fiber.Map{"message": "Nội dung landing page đã được cập nhật thành công"})
}
