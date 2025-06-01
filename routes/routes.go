package routes

import (
	"github.com/gofiber/fiber/v2"
	"notes-api/handlers"
	"notes-api/middleware"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler()
	notesHandler := handlers.NewNotesHandler()

	// API version 1 group
	api := app.Group("/api/v1")

	// Authentication routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected routes (require JWT authentication)
	protected := api.Group("", middleware.JWTMiddleware())

	// User profile route
	protected.Get("/profile", authHandler.Profile)

	// Notes routes (all protected)
	notes := protected.Group("/notes")
	notes.Post("/", notesHandler.CreateNote)           // POST /api/v1/notes
	notes.Get("/", notesHandler.GetNotes)              // GET /api/v1/notes
	notes.Get("/:id", notesHandler.GetNote)            // GET /api/v1/notes/:id
	notes.Put("/:id", notesHandler.UpdateNote)         // PUT /api/v1/notes/:id
	notes.Delete("/:id", notesHandler.DeleteNote)      // DELETE /api/v1/notes/:id
}
