package routes

import (
	"backend/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App, handler *handlers.WheelHandler, frontendURL string) {
	// Apply global logging
	app.Use(logger.New())
	
	// Apply CORS configurations
	app.Use(cors.New(cors.Config{
		AllowOrigins:     frontendURL,
		AllowHeaders:     "Origin, Content-Type, Accept, X-Edit-Token",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// API Version 1 Group
	api := app.Group("/api/v1")

	// Core check endpoint
	api.Get("/health", handler.HealthCheck)

	// Wheel operations
	api.Post("/wheels", handler.CreateWheel)
	api.Get("/wheels/:shareCode", handler.GetWheel)
	api.Put("/wheels/:id", handler.UpdateWheel)
	api.Delete("/wheels/:id", handler.DeleteWheel)

	// Spin logger
	api.Post("/wheels/:id/spin", handler.RecordSpin)

	// Log histories
	api.Get("/wheels/:id/history", handler.GetHistory)
	api.Delete("/wheels/:id/history", handler.ClearHistory)

	// Copy action
	api.Post("/wheels/:id/duplicate", handler.DuplicateWheel)

	// Analytics operations
	api.Post("/analytics/visit", handler.RecordVisit)
}
