package handler

import (
	"fib/database"
	dto "fib/dto/user"
	"fib/entity"
	"fib/repository"
	"fib/service"
	"fib/util"
	"github.com/gofiber/fiber/v2"
	"strconv"
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
		db := database.DB
		userDto := new(dto.SignupUserDataDto)

		bodyParserError := c.BodyParser(userDto)
		validationError := util.Validator.Struct(userDto)

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
		userEntity := toEnity(userDto)

		if err := db.Create(&userEntity).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				fiber.Map{
					"status":  "error",
					"message": "Couldn't create user",
					"errors":  err.Error(),
				},
			)
		}

		userResponseDto := dto.UserResponseDto{
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

func toEnity(dto *dto.SignupUserDataDto) entity.User {
	return entity.User{
		Username:  dto.Username,
		Password:  dto.Password,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
	}
}
