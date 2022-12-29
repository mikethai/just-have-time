package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	storySongRoutes "github.com/mikethai/just-have-time/internal/routes/storySong"
	userRoutes "github.com/mikethai/just-have-time/internal/routes/user"
)

func SetupRoutes(app *fiber.App) {

	root := app.Group("/", logger.New())

	// Setup the Story Song Routes
	storySongRoutes.SetupStorySongRoutes(root)

	// Setup the User Routes
	userRoutes.SetupUserRoutes(root)
}
