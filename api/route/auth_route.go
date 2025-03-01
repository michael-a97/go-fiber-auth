package route

import (
	"fib/api/handler"
	"fib/pkg/service"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(
	app fiber.Router,
	authService service.AuthService,
	userService service.UserService,
) {
	auth := app.Group("/auth")
	auth.Post("/login", handler.Login(authService, userService))
}
