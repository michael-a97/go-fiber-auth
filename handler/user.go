package handler

import (
	"fib/database"
	"fib/entity"
	"fib/util"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validToken(t *jwt.Token, id string) bool {
	n, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	claims := t.Claims.(jwt.MapClaims)

	uid := int(claims["user_id"].(float64))

	return uid == n
}

func validUser(id string, p string) bool {
	db := database.DB
	var user entity.User
	db.First(&user, id)
	if user.Username == "" {
		return false
	}
	if !CheckPasswordHash(p, user.Username) {
		return false
	}
	return true
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var user entity.User
	db.Find(&user, id)
	if user.Username == "" {
		return c.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No user found with ID",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "sucess",
			"message": "User found",
			"data":    user,
		},
	)
}

type CreateUserDto struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type UserResponseDto struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func Createuser(c *fiber.Ctx) error {
	db := database.DB
	userDto := new(CreateUserDto)

	if err := c.BodyParser(userDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Review your input",
				"errors":  err.Error(),
			},
		)
	}

	if err := util.Validator.Struct(userDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Review your input",
				"errors":  err.Error(),
			},
		)
	}

	hash, err := hashPassword(userDto.Password)
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

	userEntity := userDto.toEnity()

	if err := db.Create(&userEntity).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Couldn't create user",
				"errors":  err.Error(),
			},
		)
	}

	newUser := CreateUserDto{
		Username:  userDto.Username,
		Password:  userDto.Password,
		FirstName: userDto.FirstName,
		LastName:  userDto.LastName,
	}
	return c.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"status":  "success",
			"message": "Successfully created new user",
			"data":    newUser,
		},
	)
}

func (dto CreateUserDto) toEnity() entity.User {
	return entity.User{
		Username:  dto.Username,
		Password:  dto.Password,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
	}
}
