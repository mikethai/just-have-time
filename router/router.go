package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	searchRoutes "github.com/mikethai/just-have-time/internal/routes/search"
	storySongRoutes "github.com/mikethai/just-have-time/internal/routes/storySong"
	userRoutes "github.com/mikethai/just-have-time/internal/routes/user"
)

func SetupRoutes(app *fiber.App) {

	root := app.Group("/", logger.New())

	// Setup the Search Routes
	searchRoutes.SetupSearchRoutes(root)

	// Setup the Story Song Routes
	storySongRoutes.SetupStorySongRoutes(root)

	// Setup the User Routes
	userRoutes.SetupUserRoutes(root)
}
