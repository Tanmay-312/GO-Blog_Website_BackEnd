package main

import (
	"blog-website/database"
	"blog-website/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()

	app := fiber.New()
	routes.Setup(app)
	app.Listen(":3000")
}
