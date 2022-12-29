package userRoutes

import (
	"github.com/gofiber/fiber/v2"
	userHandler "github.com/mikethai/just-have-time/internal/handlers/user"
)

func SetupUserRoutes(router fiber.Router) {
	newsUserHandler := userHandler.NewHandler()

	storySong := router.Group("/follow")

	// Create a follow record
	storySong.Post("/", newsUserHandler.CreateUserFollow)
}
