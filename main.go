package main

import (
	"fib/database"
	"fib/api/middleware"
	"fib/pkg/repository"
	"fib/api/route"
	"fib/pkg/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
)

func main() {
	db := database.ConnectDB()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	authService := service.NewAuthService()

	app := fiber.New(
		fiber.Config{
			Prefork:       true,
			CaseSensitive: true,
			StrictRouting: true,
			ServerHeader:  "Fiber",
			AppName:       "go-fiber-auth",
		},
	)

	api := app.Group("/api", logger.New())
	route.SetupUserRouter(api, userService)
	route.SetupAuthRoutes(api, authService, userService)
	api.Get("/ping", middleware.Protected(), func(c *fiber.Ctx) error {
		return c.SendString("Pong")
	})

	log.Fatal(app.Listen(":3000"))
}
