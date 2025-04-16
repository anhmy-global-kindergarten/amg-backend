package auth

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthHandler struct {
	Router fiber.Router
	DB     *gorm.DB
}

func RegisterAuthHandler(router fiber.Router, db *gorm.DB) {
	authHandler := AuthHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Post("/register", authHandler.Register)

}
