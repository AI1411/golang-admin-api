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
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Status    string    `json:"status"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
