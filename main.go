package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mikethai/just-have-time/database"
)

func main() {
	app := fiber.New()

	// Connect to the Database
	database.ConnectDB()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Test Just have time 👋!")
	})

	app.Listen(":80")
}
