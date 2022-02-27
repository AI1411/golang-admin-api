package models

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

const DefaultPasswordCost = 14

type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	FirstName string    `json:"first_name" binding:"required,max=16"`
	LastName  string    `json:"last_name" binding:"required,max=16"`
	Image     string    `json:"image" binding:"required,max=16"`
	Age       uint8     `json:"age" binding:"required,max=16"`
	Email     string    `json:"email" binding:"required,max=16"`
	Password  []byte    `json:"password" binding:"required,max=16"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Todos     []Todo    `json:"todos" binding:"omitempty"`
}

func (user *User) SetPassword(password string) {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), DefaultPasswordCost)

	user.Password = hashPassword
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}
