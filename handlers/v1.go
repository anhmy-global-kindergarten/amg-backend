package feature

import (
	"amg-backend/cronjobs"
	"amg-backend/handlers/auth"
	"amg-backend/handlers/candidate"
	"amg-backend/handlers/comment"
	"amg-backend/handlers/landing_page"
	"amg-backend/handlers/post"
	"amg-backend/handlers/uploaded_image"
	"amg-backend/handlers/user"
	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

// @host localhost:3030
// @BasePath /amg/v1
func RegisterHandlerV1(db *mongo.Client) *fiber.App {
	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(1).Day().At("02:00").Do(func() {
		cronjobs.RunImageCleanupJob(db)
	})
	if err != nil {
		log.Fatalf("Could not schedule cron job: %v", err)
	}

	s.StartAsync()
	log.Println("Cron job scheduler started.")
	router := fiber.New()
	router.Static("/uploads", "./uploads")
	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		//AllowCredentials: true,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	v1 := router.Group("/amg/v1")
	v1.Get("/swagger/*", swagger.HandlerDefault)
	auth.RegisterAuthHandler(v1.Group("/auth-self"), db)
	user.RegisterUserHandler(v1.Group("/users"), db)
	post.RegisterPostHandler(v1.Group("/posts"), db)
	candidate.RegisterCandidateHandler(v1.Group("/candidates"), db)
	uploaded_image.RegisterUploadedImageHandler(v1.Group("/images"), db)
	landing_page.RegisterLandingPageHandler(v1.Group("/landing-page"), db)
	comment.RegisterCommentHandler(v1.Group("/comments"), db)
	return router
}
