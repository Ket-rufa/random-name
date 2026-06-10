package main

import (
	"log"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/repositories"
	"backend/internal/routes"
	"backend/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	log.Println("Starting Random Name Wheel server...")

	// 1. Parse configuration
	cfg := config.LoadConfig()

	// 2. Establish PostgreSQL GORM client
	db := database.ConnectDB(cfg.DatabaseURL)

	// 3. Instantiate dependency injection layers
	repo := repositories.NewWheelRepository(db)
	service := services.NewWheelService(repo)
	handler := handlers.NewWheelHandler(service)

	// 4. Instantiate Fiber application with centralized error handling
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	// Add panic recovery interceptor
	app.Use(recover.New())

	// 5. Setup route mappings
	routes.SetupRoutes(app, handler, cfg.FrontendURL)

	// 6. Boot HTTP listener
	log.Printf("Server boot successful, listening on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Critical: Server failed to start: %v", err)
	}
}
