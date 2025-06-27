package landing_page

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type LandingPageHandler struct {
	Router fiber.Router
	DB     *mongo.Client
}

func RegisterLandingPageHandler(router fiber.Router, db *mongo.Client) {
	postHandler := LandingPageHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Get("/get-content", postHandler.GetLandingPageContent)
	router.Post("/update-content", postHandler.UpdateLandingPageContent)
}
