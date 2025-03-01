package handler

import (
	"fib/api/presenter"
	"fib/pkg/entity"
	"fib/pkg/repository"
	"fib/pkg/service"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func GetUser(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{
					"status":  "error",
					"message": "No user found with ID",
					"data":    nil,
				},
			)
		}
		user, err := userService.FindUserById(id)

		if err == repository.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"status":  "error",
					"message": "No user found with ID",
					"data":    nil,
				},
			)
		} else {
			return c.JSON(
				fiber.Map{
					"status":  "sucess",
					"message": "User found",
					"data":    user,
				},
			)
		}

	}
}

func CreateUser(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userDto := new(struct {
			Username  string `json:"username" validate:"required"`
			Password  string `json:"password" validate:"required"`
			FirstName string `json:"first_name" validate:"required"`
			LastName  string `json:"last_name" validate:"required"`
		})

		bodyParserError := c.BodyParser(userDto)
		validationError := validator.New(
			validator.WithRequiredStructEnabled(),
		).Struct(userDto)

		var err error

		if bodyParserError != nil {
			err = bodyParserError
		} else if validationError != nil {
			err = validationError
		}

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{
					"status":  "error",
					"message": "Invalid input",
					"errors":  err.Error(),
				},
			)
		}

		hash, err := userService.HashPassword(userDto.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				fiber.Map{
					"status":  "error",
					"message": "Problem hashing password",
					"errors":  err.Error(),
				},
			)
		}

		userDto.Password = hash
		userEntity := entity.User{
			Username:  userDto.Username,
			Password:  userDto.Password,
			FirstName: userDto.FirstName,
			LastName:  userDto.LastName,
		}

		if err := userService.CreateUser(userEntity); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				fiber.Map{
					"status":  "error",
					"message": "Couldn't create user",
					"errors":  err.Error(),
				},
			)
		}

		userResponseDto := presenter.UserResponseDto{
			Username:  userDto.Username,
			FirstName: userDto.FirstName,
			LastName:  userDto.LastName,
		}

		return c.Status(fiber.StatusCreated).JSON(
			fiber.Map{
				"status":  "success",
				"message": "Successfully created new user",
				"data":    userResponseDto,
			},
		)
	}
}
