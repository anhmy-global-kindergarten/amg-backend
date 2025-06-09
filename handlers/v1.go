package feature

import (
	"amg-backend/handlers/auth"
	"amg-backend/handlers/candidate"
	"amg-backend/handlers/post"
	"amg-backend/handlers/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

// @host localhost:3030
// @BasePath /amg/v1
func RegisterHandlerV1(db *mongo.Client) *fiber.App {
	router := fiber.New()
	router.Static("/uploads", "./uploads")
	router.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	}))

	v1 := router.Group("/amg/v1")
	v1.Get("/swagger/*", swagger.HandlerDefault)
	auth.RegisterAuthHandler(v1.Group("/auth"), db)
	user.RegisterUserHandler(v1.Group("/users"), db)
	post.RegisterPostHandler(v1.Group("/posts"), db)
	candidate.RegisterCandidateHandler(v1.Group("/candidates"), db)
	return router
}
