package models

import "github.com/google/uuid"

type Product struct {
	ID        uuid.UUID `json:"id" gorm:"primary_key;index;default:uuid_generate_v4()"`
	Name      string    `json:"name" binding:"required,max=64" gorm:"type:varchar(64);not null;index"`
	Price     float64   `json:"price" binding:"required,gte=0" gorm:"type:decimal(10,2);not null;index"`
	Remarks   string    `json:"remarks" binding:"max=255" gorm:"type:varchar(255);null"`
	Quantity  int       `json:"quantity" binding:"required,gte=0" gorm:"type:int;not null"`
	CreatedAt string    `json:"created_at" gorm:"type:varchar(64);not null"`
	UpdatedAt string    `json:"updated_at" gorm:"type:varchar(64);not null"`
}
