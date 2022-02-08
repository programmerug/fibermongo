package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/programmerug/fibermongo/route"
)

func main() {
	app := fiber.New()

	// routes
	route.UserRoute(app) // add this

	app.Listen(":2022")
}
