package repository

import (
	"errors"
	"fib/pkg/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindUserById(id int64) (*entity.User, error)
	FindUserByUsername(username string) (*entity.User, error)
	CreateUser(user entity.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

var ErrRecordNotFound = errors.New("record not found")

func (r *userRepository) FindUserById(id int64) (*entity.User, error) {
	var user entity.User
	result := r.db.Find(&user, id)
	if result.RowsAffected == 0 {
		return &user, ErrRecordNotFound
	} else {
		return &user, nil
	}
}
func (r *userRepository) CreateUser(user entity.User) error {
	err := r.db.Create(&user)
	if err != nil {
		return err.Error
	}
	return nil
}

func (r *userRepository) FindUserByUsername(username string) (*entity.User, error) {
	var user entity.User

	if err := r.db.Where(&entity.User{Username: username}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
