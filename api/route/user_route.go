package route

import (
	"fib/api/handler"
	"fib/pkg/service"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRouter(app fiber.Router, userService service.UserService) {
	user := app.Group("/user")
	user.Post("/signup", handler.CreateUser(userService))
	user.Get("/", handler.GetUser(userService))
}
