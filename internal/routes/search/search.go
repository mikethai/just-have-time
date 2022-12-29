package searchRoutes

import (
	"github.com/gofiber/fiber/v2"
	searchHandler "github.com/mikethai/just-have-time/internal/handlers/search"
)

func SetupSearchRoutes(router fiber.Router) {

	newSearchHandler := searchHandler.NewHandler()
	song := router.Group("/song")
	songSearch := song.Group("/search")

	// search song by name
	songSearch.Get("/:keyWord", newSearchHandler.GetSongByKeyword)
}
