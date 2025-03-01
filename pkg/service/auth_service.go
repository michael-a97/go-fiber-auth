package service

import (
	"fib/config"
	"fib/pkg/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignToken(user *entity.User) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

func (a *authService) SignToken(user *entity.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_name"] = user.Username
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	signedToken, err := token.SignedString([]byte(config.Config("Secret")))

	return signedToken, err
}

func (a *authService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
