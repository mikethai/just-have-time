package chartsRoutes

import (
	"github.com/gofiber/fiber/v2"
	chartsHandler "github.com/mikethai/just-have-time/internal/handlers/charts"
)

func SetupChartsRoutes(router fiber.Router) {

	newChartsHandler := chartsHandler.NewHandler()
	song := router.Group("/song")
	songSearch := song.Group("/charts")

	// search song by name
	songSearch.Get("/", newChartsHandler.GetSongCharts)
}
