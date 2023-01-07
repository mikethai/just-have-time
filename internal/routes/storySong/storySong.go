package storySongRoutes

import (
	"github.com/gofiber/fiber/v2"
	storySongHandler "github.com/mikethai/just-have-time/internal/handlers/storySong"
)

func SetupStorySongRoutes(router fiber.Router) {

	newstorySongHandler := storySongHandler.NewHandler()

	// Read all story songs
	router.Get("/story-cards", newstorySongHandler.GetStorySongs)

	// Read all story songs by msno
	router.Get("/story-cards/:msno", newstorySongHandler.GetStorySongs)

	// Read all story songs without user's card
	router.Get("/story-cards-:filterUser/:msno", newstorySongHandler.GetStorySongs)

	// Create a story song
	router.Post("/story-card", newstorySongHandler.CreateStorySongs)

}
