package models

import "time"

type Epic struct {
	ID              uint64    `json:"id"`
	IsOpen          bool      `json:"is_open" binding:"omitempty,boolean"`
	AuthorID        string    `json:"author_id" binding:"required"`
	EpicTitle       string    `json:"epic_title" binding:"required"`
	EpicDescription string    `json:"epic_description" binding:"omitempty"`
	Label           string    `json:"label" binding:"omitempty"`
	MilestoneID     string    `json:"milestone_id" binding:"omitempty"`
	AssigneeID      string    `json:"assignee_id" binding:"omitempty"`
	ProjectID       string    `json:"project_id" binding:"required"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type EpicList []Epic
