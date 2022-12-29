package storySongRoutes

import (
	"github.com/gofiber/fiber/v2"
	storySongHandler "github.com/mikethai/just-have-time/internal/handlers/storySong"
)

func SetupStorySongRoutes(router fiber.Router) {

	newstorySongHandler := storySongHandler.NewHandler()
	storySong := router.Group("/story-song")

	// Create a story song
	storySong.Post("/", newstorySongHandler.CreateStorySongs)
	// Read all story songs
	storySong.Get("/", newstorySongHandler.GetStorySongs)
	// Read one story song
	storySong.Get("/:storySongId", newstorySongHandler.GetStorySong)
	// Update one story song
	storySong.Put("/:storySongId", newstorySongHandler.UpdateStorySong)
}
