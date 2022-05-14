package models

import (
	"time"

	"github.com/google/uuid"
)

type Coupon struct {
	ID                string    `json:"id"`
	Title             string    `json:"title" binding:"required"`
	Remarks           string    `json:"remarks" binding:"omitempty,max=255"`
	DiscountAmount    uint64    `json:"discount_amount" binding:"omitempty,min=1"`
	DiscountRate      uint8     `json:"discount_rate" binding:"omitempty,min=1,max=100"`
	MaxDiscountAmount uint64    `json:"max_discount_amount" binding:"omitempty,min=1"`
	UseStartAt        time.Time `json:"use_start_at" binding:"required"`
	UseEndAt          time.Time `json:"use_end_at" binding:"required"`
	PublicStartAt     time.Time `json:"public_start_at" binding:"required"`
	PublicEndAt       time.Time `json:"public_end_at" binding:"required"`
	IsPublic          bool      `json:"is_public"`
	IsPremium         bool      `json:"is_premium"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (c *Coupon) CreateUUID() {
	newUUID, _ := uuid.NewRandom()
	c.ID = newUUID.String()
}
