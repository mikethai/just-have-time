package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Test Just have time 👋!")
	})

	app.Listen(":80")
}
