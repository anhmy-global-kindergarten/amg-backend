package comment

import (
	"amg-backend/config"
	"amg-backend/models"
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// GetCommentsByPostId godoc
// @Summary Get comments by post-ID
// @Description Get all comments for a specific post
// @Tags Comment
// @Accept json
// @Produce json
// @Param postId query string true "Post ID"
// @Success 200 {array} models.Comment
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/comments/get-comments-by-post [get]
func (h *CommentHandler) GetCommentsByPostId(c *fiber.Ctx) error {
	postId := c.Query("postId")
	if postId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "postId is required"})
	}

	var comments []models.Comment
	collection := h.DB.Database(config.DBName).Collection("Comment")

	filter := bson.M{
		"post_id": postId,
		"status":  bson.M{"$ne": "deleted"},
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &comments); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode comments"})
	}

	if comments == nil {
		comments = make([]models.Comment, 0)
	}

	return c.JSON(comments)
}

// CreateComment godoc
// @Summary Create a new comment
// @Description Create a new comment for a post
// @Tags Comment
// @Accept json
// @Produce json
// @Param comment body models.CreateCommentPayload true "Data to create a comment"
// @Success 201 {object} models.Comment
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/comments/create-comment [post]
func (h *CommentHandler) CreateComment(c *fiber.Ctx) error {
	var payload models.CreateCommentPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if payload.PostId == "" || payload.AuthorName == "" || payload.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "postId, authorName, và content là bắt buộc"})
	}

	comment := models.Comment{
		ID:         primitive.NewObjectID(),
		PostId:     payload.PostId,
		AuthorId:   payload.AuthorId,
		AuthorName: payload.AuthorName,
		Content:    payload.Content,
		Status:     "new",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	collection := h.DB.Database(config.DBName).Collection("Comment")
	_, err := collection.InsertOne(context.TODO(), comment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create comment"})
	}

	return c.Status(fiber.StatusCreated).JSON(comment)
}

// DeleteComment godoc
// @Summary Delete a comment
// @Description Mark a comment as deleted by ID
// @Tags Comment
// @Accept json
// @Produce json
// @Param id path string true "Comment ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/comments/delete-comment/{id} [post]
func (h *CommentHandler) DeleteComment(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID is required"})
	}

	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	updateData := bson.M{
		"updated_at": time.Now(),
		"status":     "deleted",
	}

	collection := h.DB.Database(config.DBName).Collection("Comment")
	result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": updateData})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Delete failed"})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Comment not found"})
	}

	return c.JSON(fiber.Map{"message": "Comment deleted successfully"})
}

// UpdateComment godoc
// @Summary Update a comment
// @Description Update a comment by ID
// @Tags Comment
// @Accept json
// @Produce json
// @Param id path string true "Comment ID"
// @Param comment body models.Comment true "Comment data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/comments/update-comment/{id} [post]
func (h *CommentHandler) UpdateComment(c *fiber.Ctx) error {
	var comment models.Comment
	if err := c.BodyParser(&comment); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post ID format"})
	}

	updateData := bson.M{
		"updated_at": time.Now(),
		"content":    comment.Content,
	}

	collection := h.DB.Database(config.DBName).Collection("Comment")
	_, err = collection.UpdateByID(context.TODO(), id, bson.M{"$set": updateData})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB error"})
	}
	return c.JSON(fiber.Map{"message": "Comment updated successfully"})
}
