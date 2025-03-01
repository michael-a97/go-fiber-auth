package main

import (
	"fib/database"
	"fib/handler"
	"fib/middleware"
	"fib/repository"
	"fib/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
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
			AppName:       "App name",
		},
	)

	api := app.Group("/api", logger.New())
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login(authService, userService))
	user := api.Group("/user")
	user.Post("/signup", handler.CreateUser(userService))
	user.Get("/", handler.GetUser(userService))
	api.Get("/ping", middleware.Protected(), func(c *fiber.Ctx) error {
		return c.SendString("Pong")
	})



	log.Fatal(app.Listen(":3000"))
}
