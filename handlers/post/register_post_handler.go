package post

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostHandler struct {
	Router fiber.Router
	DB     *mongo.Client
}

func RegisterPostHandler(router fiber.Router, db *mongo.Client) {
	postHandler := PostHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Get("/get-all-posts", postHandler.GetAllPosts)
	router.Get("/get-post/:id", postHandler.GetPostById)
	router.Get("/get-posts-by-category/:category", postHandler.GetPostsByCategory)
	router.Get("/get-single-post-by-category/:category", postHandler.GetSinglePostByCategory)
	router.Get("/get-posts-by-status/:status", postHandler.GetPostsByStatus)
	router.Post("/update-post/:id", postHandler.UpdatePost)
	router.Post("/create-post", postHandler.CreatePost)
	router.Post("/delete-post/:id", postHandler.DeletePost)
	router.Post("/recovery-post/:id", postHandler.RecoveryPost)
}
