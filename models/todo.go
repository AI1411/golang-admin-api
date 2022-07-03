package models

import "time"

type TodoStatus string

const (
	TodoStatusSuccess    TodoStatus = "success"
	TodoStatusWaiting    TodoStatus = "waiting"
	TodoStatusCanceled   TodoStatus = "canceled"
	TodoStatusProcessing TodoStatus = "processing"
	TodoStatusDone       TodoStatus = "done"
)

type Todo struct {
	ID        uint64    `json:"id" gorm:"primaryKey" example:"1"`
	Title     string    `json:"title" binding:"required,max=64" example:"タイトル"`
	Body      string    `json:"body" binding:"required,max=64" example:"本文"`
	Status    string    `json:"status" binding:"required,oneof=new processing done closed" example:"success"`
	UserID    *string   `json:"user_id" binding:"required" example:"14841545-8a11-47d1-bf95-59a1c7f1d8ec"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
