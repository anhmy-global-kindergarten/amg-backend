package candidate

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

// GetAllCandidates godoc
// @Summary Get all candidate
// @Description Retrieves all candidates from the database
// @Tags candidate
// @Accept json
// @Produce json
// @Success 200 {array} models.Candidate
// @Failure 500 {object} map[string]string
// @Router /amg/v1/candidates/get-all-candidates [get]
func (h *CandidateHandler) GetAllCandidates(c *fiber.Ctx) error {
	var candidates []models.Candidate
	collection := h.DB.Database(config.DBName).Collection("Candidate")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB error"})
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var candidate models.Candidate
		cursor.Decode(&candidate)
		candidates = append(candidates, candidate)
	}
	return c.JSON(candidates)
}

// GetCandidateById godoc
// @Summary Get candidate by ID
// @Description Retrieves a candidate by its ID
// @Tags candidate
// @Accept json
// @Produce json
// @Param id path string true "Candidate ID"
// @Success 200 {object} models.Candidate
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/candidates/get-candidate/{id} [get]
func (h *CandidateHandler) GetCandidateById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid candidate ID"})
	}

	var candidate models.Candidate
	collection := h.DB.Database(config.DBName).Collection("Candidate")

	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&candidate)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Candidate not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}

	return c.JSON(candidate)
}

// GetCandidatesByStatus godoc
// @Summary Get candidates by status
// @Description Retrieves candidates by their status
// @Tags candidate
// @Accept json
// @Produce json
// @Param status path string true "Candidate Status"
// @Success 200 {array} models.Candidate
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/candidates/get-candidates-by-status/{status} [get]
func (h *CandidateHandler) GetCandidatesByStatus(c *fiber.Ctx) error {
	status := c.Params("status")
	if status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid status"})
	}

	var candidates []models.Candidate
	collection := h.DB.Database(config.DBName).Collection("Post")
	cursor, err := collection.Find(context.TODO(), bson.M{"status": status})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}
	// Tìm theo status
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var candidate models.Candidate
		if err := cursor.Decode(&candidate); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Decode error"})
		}
		candidates = append(candidates, candidate)
	}

	return c.JSON(candidates)
}

// UpdateCandidate godoc
// @Summary Update a candidate
// @Description Updates a candidate's information by its ID
// @Tags candidate
// @Accept json
// @Produce json
// @Param id path string true "Candidate ID"
// @Param candidate body models.Candidate true "Candidate data to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/candidates/update-candidate/{id} [post]
func (h *CandidateHandler) UpdateCandidate(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, _ := primitive.ObjectIDFromHex(idParam)

	updateData := bson.M{}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}
	updateData["update_at"] = time.Now()

	collection := h.DB.Database(config.DBName).Collection("Candidate")
	_, err := collection.UpdateByID(context.TODO(), id, bson.M{"$set": updateData})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "update failed"})
	}
	return c.JSON(fiber.Map{"message": "updated"})
}

// CreateCandidate godoc
// @Summary Create a new candidate
// @Description Creates a new candidate in the database
// @Tags candidate
// @Accept json
// @Produce json
// @Param candidate body models.Candidate true "Candidate data"
// @Success 200 {object} models.Candidate
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/candidates/create-candidate [post]
func (h *CandidateHandler) CreateCandidate(c *fiber.Ctx) error {
	var candidate models.Candidate
	if err := c.BodyParser(&candidate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	candidate.ID = primitive.NewObjectID()
	candidate.CreateAt = time.Now()
	candidate.UpdateAt = time.Now()
	candidate.Status = "new" // Default status

	collection := h.DB.Database(config.DBName).Collection("Candidate")
	_, err := collection.InsertOne(context.TODO(), candidate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create candidate"})
	}

	return c.JSON(candidate)
}

// DeleteCandidate godoc
// @Summary Delete a candidate
// @Description Deletes a candidate by its ID
// @Tags candidate
// @Accept json
// @Produce json
// @Param id path string true "Candidate ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/candidates/delete-candidate/{id} [post]
func (h *CandidateHandler) DeleteCandidate(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, _ := primitive.ObjectIDFromHex(idParam)

	updateData := bson.M{}
	updateData["update_at"] = time.Now()
	updateData["status"] = "deleted"

	collection := h.DB.Database(config.DBName).Collection("Candidate")
	_, err := collection.UpdateByID(context.TODO(), id, bson.M{"$set": updateData})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "delete failed"})
	}
	return c.JSON(fiber.Map{"message": "deleted"})
}

// RecoveryCandidate godoc
// @Summary Recover a deleted candidate
// @Description Recovers a candidate by its ID
// @Tags candidate
// @Accept json
// @Produce json
// @Param id path string true "Candidate ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/candidates/recovery-candidate/{id} [post]
func (h *CandidateHandler) RecoveryCandidate(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, _ := primitive.ObjectIDFromHex(idParam)

	updateData := bson.M{}
	updateData["update_at"] = time.Now()
	updateData["status"] = "recovered"

	collection := h.DB.Database(config.DBName).Collection("Candidate")
	_, err := collection.UpdateByID(context.TODO(), id, bson.M{"$set": updateData})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "recovery failed"})
	}
	return c.JSON(fiber.Map{"message": "recovered"})
}
