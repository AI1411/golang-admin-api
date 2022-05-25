package models

import (
	"github.com/google/uuid"
	"time"
)

type UserGroup struct {
	Id        string    `json:"id"`
	GroupName string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *UserGroup) CreateUUID() string {
	newUUID, _ := uuid.NewRandom()
	return newUUID.String()
}
