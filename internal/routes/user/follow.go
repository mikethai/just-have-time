package userRoutes

import (
	"github.com/gofiber/fiber/v2"
	userHandler "github.com/mikethai/just-have-time/internal/handlers/user"
)

func SetupUserRoutes(router fiber.Router) {
	newsUserHandler := userHandler.NewHandler()

	storySong := router.Group("/follow")

	storySong.Post("/sync", newsUserHandler.SyncUserFollow)
}
