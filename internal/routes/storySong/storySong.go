package storySongRoutes

import (
	"github.com/gofiber/fiber/v2"
	storySongHandler "github.com/mikethai/just-have-time/internal/handlers/storySong"
)

func SetupStorySongRoutes(router fiber.Router) {

	newstorySongHandler := storySongHandler.NewHandler()

	// Read all story songs
	router.Get("/story-cards", newstorySongHandler.GetStorySongs)

	// Create a story song
	router.Post("/story-card", newstorySongHandler.CreateStorySongs)

	// storySong := router.Group("/story-card")
	// Read one story song
	// storySong.Get("/:storySongId", newstorySongHandler.GetStorySong)
	// Update one story song
	// storySong.Put("/:storySongId", newstorySongHandler.UpdateStorySong)
}
