package comment

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentHandler struct {
	Router fiber.Router
	DB     *mongo.Client
}

func RegisterCommentHandler(router fiber.Router, db *mongo.Client) {
	commentHandler := CommentHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Get("/get-comments-in-post", commentHandler.GetCommentsByPostId)
	router.Post("/update-comment/:id", commentHandler.UpdateComment)
	router.Post("/create-comment", commentHandler.CreateComment)
	router.Post("/delete-comment/:id", commentHandler.DeleteComment)
}
