package models

import (
	"github.com/google/uuid"
)

type Project struct {
	ID                 string `json:"id"`
	ProjectTitle       string `json:"project_title" binding:"required,max=64"`
	ProjectDescription string `json:"project_description" binding:"omitempty,max=255"`

	Epics EpicList `json:"epics" binding:"omitempty,dive"`
}

func (p *Project) CreateUUID() {
	newUUID, _ := uuid.NewRandom()
	p.ID = newUUID.String()
}
