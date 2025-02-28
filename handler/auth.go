package handler

import (
	"errors"
	"fib/config"
	"fib/database"
	"fib/entity"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserByUsername(username string) (*entity.User, error) {
	db := database.DB
	var user entity.User

	fmt.Printf("db=%v dbt=%T", db, db)
	if err := db.Where(&entity.User{Username: username}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	input := new(LoginInput)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"status": "error",
				"message": "Invalid username or password",
				"errors":  err.Error(),
			},
		)
	}

	username := input.Username
	password := input.Password

	err := *new(error)
	var user *entity.User
	fmt.Printf("username=%s\n", username)
	fmt.Printf("password=%s\n", password)
	user, err = getUserByUsername(username)

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
			fiber.Map{"status": "error",
				"message": "Invalid username or password",
				"data":    err,
			},
		)
	}

	if !CheckPasswordHash(password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Invalid username or password",
				"data":    nil,
			},
		)
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_name"] = user.Username
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	signedToken, err := token.SignedString([]byte(config.Config("Secret")))

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
