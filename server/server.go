package server

import (
	"amg-backend/config"
	"amg-backend/database"
	h "amg-backend/handlers"
	"log"
)

func RegisterServer() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}

	// Initialize router
	r := h.RegisterHandlerV1(db)
	// Start server
	log.Printf("Server is running on port %s", cfg.ServerPort)
	if err := r.Listen(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
