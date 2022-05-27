package models

import (
	"github.com/google/uuid"
	"time"
)

type Milestone struct {
	ID                   string    `json:"id"`
	MilestoneTitle       string    `json:"milestone_title" binding:"required,max=64"`
	MilestoneDescription string    `json:"milestone_description" binding:"omitempty,max=255"`
	ProjectID            string    `json:"project_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func (m *Milestone) CreateUUID() {
	newUUID, _ := uuid.NewRandom()
	m.ID = newUUID.String()
}
