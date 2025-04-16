package feature

import (
	"amg-backend/handlers/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
)

// @host localhost:8080
// @BasePath /dbms/v1
func RegisterHandlerV1(db *gorm.DB) *fiber.App {
	router := fiber.New()
	v1 := router.Group("/livenest/v1")
	v1.Get("/swagger/*", swagger.HandlerDefault)
	auth.RegisterAuthHandler(v1.Group("/auth"), db)
	return router
}
