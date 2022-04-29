package models

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID          uuid.UUID `json:"id"`
	ProductName string    `json:"product_name" binding:"required,max=64"`
	Price       uint      `json:"price" binding:"required,gte=0"`
	Remarks     string    `json:"remarks" binding:"omitempty,max=255"`
	Quantity    int       `json:"quantity" binding:"required,gte=0"`
	CreatedAt   time.Time `json:"created_at" binding:"omitempty"`
	UpdatedAt   time.Time `json:"updated_at" binding:"omitempty"`
}

func (p *Product) CreateUUID() {
	p.ID = uuid.New()
}
