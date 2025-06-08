package feature

import (
	"amg-backend/handlers/auth"
	"amg-backend/handlers/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

// @host localhost:3030
// @BasePath /amg/v1
func RegisterHandlerV1(db *mongo.Client) *fiber.App {
	router := fiber.New()
	v1 := router.Group("/amg/v1")
	v1.Get("/swagger/*", swagger.HandlerDefault)
	auth.RegisterAuthHandler(v1.Group("/auth"), db)
	user.RegisterUserHandler(v1.Group("/user"), db)
	return router
}
