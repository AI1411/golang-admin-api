package models

import "time"

type Issue struct {
	ID          string    `json:"id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"omitempty,max=255"`
	UserID      string    `json:"user_id" binding:"omitempty"`
	MilestoneID string    `json:"milestone_id" binding:"omitempty"`
	IssueStatus string    `json:"issue_status" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Milestone *Milestone `json:"milestone"`
}
