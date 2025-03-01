package handler

import (
	"fib/pkg/entity"
	"fib/pkg/service"
	"github.com/gofiber/fiber/v2"
)

func Login(authService service.AuthService, userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type LoginInput struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		input := new(LoginInput)

		if err := c.BodyParser(input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{
					"status":  "error",
					"message": "Invalid username or password",
					"errors":  err.Error(),
				},
			)
		}

		username := input.Username
		password := input.Password

		err := *new(error)
		var user *entity.User
		user, err = userService.FindUserByUsername(username)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				fiber.Map{
					"status":  "error",
					"message": "Internal server error",
					"data":    err,
				},
			)
		} else if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				fiber.Map{
					"status":  "error",
					"message": "Invalid username or password",
					"data":    err,
				},
			)
		}

		if !authService.CheckPasswordHash(password, user.Password) {
			return c.Status(fiber.StatusUnauthorized).JSON(
				fiber.Map{
					"status":  "error",
					"message": "Invalid username or password",
					"data":    nil,
				},
			)
		}

		signedToken, err := authService.SignToken(user)

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(
			fiber.Map{
				"status":  "Success",
				"message": "Successfully logged in",
				"data":    signedToken,
			},
		)
	}

}
