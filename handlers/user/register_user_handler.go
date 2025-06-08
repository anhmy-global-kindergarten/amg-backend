package user

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	Router fiber.Router
	DB     *mongo.Client
}

func RegisterUserHandler(router fiber.Router, db *mongo.Client) {
	userHandler := UserHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Get("/get-all-user", userHandler.GetAllUsers)
	router.Get("/get-user/:id", userHandler.GetAllUsers)
	router.Post("/update-user/:id", userHandler.UpdateUser)
	router.Post("/deactivate-user/:id", userHandler.DeactivateUser)
	router.Post("/reactivate-user/:id", userHandler.ReactivateUser)
}
