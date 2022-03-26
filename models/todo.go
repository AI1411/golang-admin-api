package models

import "time"

type TodoStatus string

const (
	TodoStatusSuccess    TodoStatus = "success"
	TodoStatusCanceled   TodoStatus = "canceled"
	TodoStatusProcessing TodoStatus = "processing"
	TodoStatusDone       TodoStatus = "done"
)

type Todo struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" binding:"required,max=64"`
	Body      string    `json:"body" binding:"required,max=64"`
	Status    string    `json:"status" binding:"required,max=64"`
	UserID    string    `json:"user_id" binding:"required,max=64"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
