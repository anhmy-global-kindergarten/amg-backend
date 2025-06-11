package post

import (
	"amg-backend/config"
	"amg-backend/models"
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"path/filepath"
	"regexp"
	"time"
)

// GetAllPosts godoc
// @Summary Get all posts
// @Description Retrieves all posts from the database
// @Tags post
// @Accept json
// @Produce json
// @Success 200 {array} models.Post
// @Failure 500 {object} map[string]string
// @Router /amg/v1/posts/get-all-posts [get]
func (h *PostHandler) GetAllPosts(c *fiber.Ctx) error {
	var posts []models.Post
	collection := h.DB.Database(config.DBName).Collection("Post")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB error"})
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var post models.Post
		cursor.Decode(&post)
		posts = append(posts, post)
	}
	return c.JSON(posts)
}

// GetPostById godoc
// @Summary Get a single post by ID with associated images
// @Description Retrieves a post by its ID
// @Tags post
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} models.PostDetailResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/posts/get-post/{id} [get]
func (h *PostHandler) GetPostById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post ID"})
	}

	var post models.Post
	collection := h.DB.Database(config.DBName).Collection("Post")

	// Tìm theo _id
	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&post)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}

	re := regexp.MustCompile(regexp.QuoteMeta(config.BaseURL) + `/uploads/[^"]+`)
	imageUrls := re.FindAllString(post.Content, -1)

	var relatedImages []models.UploadedImage
	if len(imageUrls) > 0 {
		imageCollection := h.DB.Database(config.DBName).Collection("UploadedImage")
		filter := bson.M{"url": bson.M{"$in": imageUrls}}

		cursor, err := imageCollection.Find(context.TODO(), filter)
		if err == nil {
			defer cursor.Close(context.TODO())
			cursor.All(context.TODO(), &relatedImages)
		}
	}

	response := models.PostDetailResponse{
		Post:   post,
		Images: relatedImages,
	}

	return c.JSON(response)
}

// GetPostsByCategory godoc
// @Summary Get posts by category
// @Description Retrieves posts by category
// @Tags post
// @Accept json
// @Produce json
// @Param category path string true "Post Category"
// @Success 200 {array} models.Post
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/posts/get-posts-by-category/{category} [get]
func (h *PostHandler) GetPostsByCategory(c *fiber.Ctx) error {
	category := c.Params("category")
	if category == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category"})
	}

	var posts []models.Post
	collection := h.DB.Database(config.DBName).Collection("Post")
	cursor, err := collection.Find(context.TODO(), bson.M{"category": category})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}
	// Tìm theo category
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Decode error"})
		}
		posts = append(posts, post)
	}

	return c.JSON(posts)
}

// GetPostsByStatus godoc
// @Summary Get posts by status
// @Description Retrieves posts by status
// @Tags post
// @Accept json
// @Produce json
// @Param status path string true "Post Status"
// @Success 200 {array} models.Post
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/posts/get-posts-by-status/{status} [get]
func (h *PostHandler) GetPostsByStatus(c *fiber.Ctx) error {
	status := c.Params("status")
	if status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid status"})
	}

	var posts []models.Post
	collection := h.DB.Database(config.DBName).Collection("Post")
	cursor, err := collection.Find(context.TODO(), bson.M{"status": status})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}
	// Tìm theo status
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Decode error"})
		}
		posts = append(posts, post)
	}

	return c.JSON(posts)
}

// UpdatePost godoc
// @Summary Update a post
// @Description Updates a post by its ID
// @Tags post
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param body body models.Post true "Post data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/posts/update-post/{id} [post]
func (h *PostHandler) UpdatePost(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, _ := primitive.ObjectIDFromHex(idParam)

	var updateData bson.M
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}
	updateData["update_at"] = time.Now()

	collection := h.DB.Database(config.DBName).Collection("Post")
	_, err := collection.UpdateByID(context.TODO(), id, bson.M{"$set": updateData})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "update failed"})
	}
	return c.JSON(fiber.Map{"message": "updated"})
}

// CreatePost godoc
// @Summary Create a new post
// @Description Creates a new post. The content should contain full URLs to images previously uploaded.
// @Tags post
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Post Title"
// @Param content formData string true "Post Content"
// @Param category formData string true "Post Category"
// @Param author formData string true "Post Author"
// @Param headerImage formData file false "Header Image"
// @Success 200 {object} models.Post
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/posts/create-post [post]
func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse form"})
	}

	title := form.Value["title"][0]
	content := form.Value["content"][0]
	category := form.Value["category"][0]
	author := form.Value["author"][0]

	var headerImagePath string
	file, err := c.FormFile("headerImage")
	if err == nil && file != nil {
		// LƯU Ý: Nên dùng tên duy nhất cho header image để tránh ghi đè
		uniqueFilename := uuid.New().String() + filepath.Ext(file.Filename)
		savePath := fmt.Sprintf("./uploads/%s", uniqueFilename)
		if err := c.SaveFile(file, savePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save header image"})
		}
		headerImagePath = fmt.Sprintf("%s/uploads/%s", config.BaseURL, uniqueFilename)
	}

	//if err := service.MarkImagesAsUsed(h.DB, content); err != nil {
	//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": `Xảy ra lỗi khi đánh dấu ảnh đã sử dụng: ${err}`})
	//}

	cleanContent := cleanImageStyles(content)

	post := models.Post{
		ID:          primitive.NewObjectID(),
		Title:       title,
		Content:     cleanContent,
		Category:    category,
		Author:      author,
		HeaderImage: headerImagePath,
		CreateAt:    time.Now(),
		UpdateAt:    time.Now(),
		Status:      "active",
	}

	collection := h.DB.Database(config.DBName).Collection("Post")
	_, err = collection.InsertOne(context.TODO(), &post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create post"})
	}

	return c.JSON(post)
}

func cleanImageStyles(htmlContent string) string {
	re := regexp.MustCompile(`<img([^>]+)style="[^"]*"([^>]*)>`)
	return re.ReplaceAllString(htmlContent, `<img$1$2>`)
}

// DeletePost godoc
// @Summary Delete a post
// @Description Deletes a post by its ID
// @Tags post
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/posts/delete-post/{id} [post]
func (h *PostHandler) DeletePost(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, _ := primitive.ObjectIDFromHex(idParam)

	var updateData bson.M
	updateData["update_at"] = time.Now()
	updateData["status"] = "deleted"

	collection := h.DB.Database(config.DBName).Collection("Post")
	_, err := collection.UpdateByID(context.TODO(), id, bson.M{"$set": updateData})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "delete failed"})
	}
	return c.JSON(fiber.Map{"message": "deleted"})
}

// RecoveryPost godoc
// @Summary Recover a deleted post
// @Description Recovers a post by its ID
// @Tags post
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /amg/v1/posts/recovery-post/{id} [post]
func (h *PostHandler) RecoveryPost(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, _ := primitive.ObjectIDFromHex(idParam)

	var updateData bson.M
	updateData["update_at"] = time.Now()
	updateData["status"] = "active"

	collection := h.DB.Database(config.DBName).Collection("Post")
	_, err := collection.UpdateByID(context.TODO(), id, bson.M{"$set": updateData})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "recovery failed"})
	}
	return c.JSON(fiber.Map{"message": "recovered"})
}
