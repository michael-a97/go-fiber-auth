package service

import (
	"fib/entity"
	"fib/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	HashPassword(password string) (string, error)
	FindUserByUsername(username string) (*entity.User, error)
	FindUserById(id int64) (*entity.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *userService) FindUserByUsername(username string) (*entity.User, error) {
	return s.userRepository.FindUserByUsername(username)
}

func (s *userService) FindUserById(id int64) (*entity.User, error) {
	return s.userRepository.FindUserById(id)
}
