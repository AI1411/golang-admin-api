package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID                 string    `json:"id"`
	ProjectTitle       string    `json:"project_title" binding:"required,max=64"`
	ProjectDescription string    `json:"project_description" binding:"omitempty,max=255"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (p *Project) CreateUUID() {
	newUUID, _ := uuid.NewRandom()
	p.ID = newUUID.String()
}
