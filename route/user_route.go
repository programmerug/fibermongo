package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/programmerug/fibermongo/controller"
)

func UserRoute(app *fiber.App) {
	app.Post("/user", controller.CreateUser)
	app.Get("/user/:userId", controller.GetAUser)
	app.Put("/user/:userId", controller.EditAUser)
	app.Delete("/user/:userId", controller.DeleteAUser)
	app.Get("/users", controller.GetAllUsers)
}
