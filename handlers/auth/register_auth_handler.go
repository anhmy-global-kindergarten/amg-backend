package auth

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	Router fiber.Router
	DB     *mongo.Client
}

func RegisterAuthHandler(router fiber.Router, db *mongo.Client) {
	authHandler := AuthHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Post("/register", authHandler.Register)
	router.Post("/login", authHandler.Login)

}
