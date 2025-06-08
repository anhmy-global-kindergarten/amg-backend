package user

import (
	"amg-backend/config"
	"amg-backend/models"
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieves all users from the database
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]string
// @Router /amg/v1/users/get-all-user [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	var users []models.User
	collection := h.DB.Database(config.DBName).Collection("User")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB error"})
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var user models.User
		cursor.Decode(&user)
		users = append(users, user)
	}
	return c.JSON(users)
}

// GetUsersById godoc
// @Summary Get user by ID
// @Description Retrieves a user by their ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/users/get-user/{id} [get]
func (h *UserHandler) GetUsersById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var user models.User
	collection := h.DB.Database(config.DBName).Collection("User")

	// Tìm theo _id
	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}

	return c.JSON(user)
}

// UpdateUser godoc
// @Summary Update user information
// @Description Updates user information based on the provided ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param body body models.User true "User data to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/users/update-user/{id} [post]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, _ := primitive.ObjectIDFromHex(idParam)

	var updateData bson.M
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}
	updateData["update_at"] = time.Now()

	collection := h.DB.Database(config.DBName).Collection("User")
	_, err := collection.UpdateByID(context.TODO(), id, bson.M{"$set": updateData})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "update failed"})
	}
	return c.JSON(fiber.Map{"message": "updated"})
}

// DeactivateUser godoc
// @Summary Deactivate a user
// @Description Deactivates a user account by setting isActive to false
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/users/deactivate-user/{id} [post]
func (h *UserHandler) DeactivateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, _ := primitive.ObjectIDFromHex(idParam)

	var updateData bson.M
	updateData["update_at"] = time.Now()
	updateData["isActive"] = false

	collection := h.DB.Database(config.DBName).Collection("User")
	_, err := collection.UpdateByID(context.TODO(), id, bson.M{"$set": updateData})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "update failed"})
	}
	return c.JSON(fiber.Map{"message": "deactivated"})
}

// ReactivateUser godoc
// @Summary Reactivate a user
// @Description Reactivates a previously deactivated user account
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/users/reactivate-user/{id} [post]
func (h *UserHandler) ReactivateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, _ := primitive.ObjectIDFromHex(idParam)

	var updateData bson.M
	updateData["update_at"] = time.Now()
	updateData["isActive"] = true

	collection := h.DB.Database(config.DBName).Collection("User")
	_, err := collection.UpdateByID(context.TODO(), id, bson.M{"$set": updateData})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "update failed"})
	}
	return c.JSON(fiber.Map{"message": "reactivated"})
}
