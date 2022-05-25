package models

import (
	"github.com/google/uuid"
	"time"
)

type UserGroup struct {
	Id        string    `json:"id"`
	GroupName string    `json:"group_name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *UserGroup) CreateUUID() {
	newUUID, _ := uuid.NewRandom()
	u.Id = newUUID.String()
}
