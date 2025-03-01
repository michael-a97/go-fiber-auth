package repository

import (
	"errors"
	"fib/database"
	"fib/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindUserById(id int64) (*entity.User, error)
	FindUserByUsername(username string) (*entity.User, error)
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
	db := database.DB
	var user entity.User
	result := db.Find(&user, id)
	if result.RowsAffected == 0 {
		return &user, ErrRecordNotFound
	} else {
		return &user, nil
	}
}

func (r *userRepository) FindUserByUsername(username string) (*entity.User, error) {
	db := database.DB
	var user entity.User

	if err := db.Where(&entity.User{Username: username}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
