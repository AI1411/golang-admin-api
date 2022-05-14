package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const DefaultPasswordCost = 14

type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name" binding:"required,max=16"`
	LastName  string    `json:"last_name" binding:"required,max=16"`
	Age       uint8     `json:"age" binding:"required,max=16"`
	Email     string    `json:"email" binding:"required,max=16"`
	Password  []byte    `json:"password" binding:"required,max=16"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Todos     []Todo    `json:"todos" binding:"omitempty"`
}

func (u *User) SetPassword(password string) {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), DefaultPasswordCost)

	u.Password = hashPassword
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(u.Password, []byte(password))
}

func (u *User) CreateUUID() {
	newUUID, _ := uuid.NewRandom()
	u.ID = newUUID.String()
}
