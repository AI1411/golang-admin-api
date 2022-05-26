package models

import (
	"github.com/google/uuid"
	"time"
)

type UserGroup struct {
	ID        string    `json:"id"`
	GroupName string    `json:"group_name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Users Users `json:"users" gorm:"many2many:group_user;"`
}

func (u *UserGroup) CreateUUID() {
	newUUID, _ := uuid.NewRandom()
	u.ID = newUUID.String()
}
