package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mikethai/just-have-time/database"
	"github.com/mikethai/just-have-time/router"
)

func main() {
	app := fiber.New()

	// Connect to the Database
	database.ConnectDB()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	// Setup the router
	router.SetupRoutes(app)

	app.Listen(":80")
}
