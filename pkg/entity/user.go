package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string `gorm:"not null" validate:"required,min=3,max=50" json:"username"`
	Password  string `gorm:"uniqueIndex;not null" validate:"required,min=6,max=50" json:"password"`
	FirstName string `gorm:"not null" validate:"required,min=3,max=100" json:"first_name"`
	LastName  string `gorm:"not null" validate:"required,min=3,max=100" json:"last_name"`
}
