package candidate

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type CandidateHandler struct {
	Router fiber.Router
	DB     *mongo.Client
}

func RegisterCandidateHandler(router fiber.Router, db *mongo.Client) {
	candidateHandler := CandidateHandler{
		Router: router,
		DB:     db,
	}

	// Register all endpoints here
	router.Get("/get-all-candidates", candidateHandler.GetAllCandidates)
	router.Get("/get-candidate/:id", candidateHandler.GetCandidateById)
	router.Get("/get-candidates-by-status/:status", candidateHandler.GetCandidatesByStatus)
	router.Post("/update-candidate/:id", candidateHandler.UpdateCandidate)
	router.Post("/create-candidate", candidateHandler.CreateCandidate)
	router.Post("/delete-candidate/:id", candidateHandler.DeleteCandidate)
	router.Post("/recovery-candidate/:id", candidateHandler.RecoveryCandidate)
}
