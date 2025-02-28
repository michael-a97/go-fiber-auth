package router

import (
	"fib/handler"
	"fib/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api",logger.New())

	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)

	user := api.Group("/user")
	user.Post("/signup", handler.Createuser)
	user.Get("/", handler.GetUser)

	api.Get("/ping", middleware.Protected(), func(c *fiber.Ctx) error {
		return c.SendString("Pong")
	})
}
