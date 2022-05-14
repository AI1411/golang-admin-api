package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusNew           OrderStatus = "new"
	OrderStatusPaid          OrderStatus = "paid"
	OrderStatusCancelled     OrderStatus = "cancelled"
	OrderStatusDelivered     OrderStatus = "delivered"
	OrderStatusRefunded      OrderStatus = "refunded"
	OrderStatusReturned      OrderStatus = "returned"
	OrderStatusPartially     OrderStatus = "partially"
	OrderStatusPartiallyPaid OrderStatus = "partially_paid"
)

type Order struct {
	ID           string          `json:"id"`
	UserID       int64           `json:"user_id" binding:"required,min=1"`
	Quantity     int64           `json:"quantity"`
	TotalPrice   int64           `json:"total_price"`
	OrderStatus  OrderStatus     `json:"order_status" binding:"omitempty,oneof=new paid cancelled delivered refunded returned partially partially_paid"`
	Remarks      string          `json:"remarks" binding:"omitempty,max=255"`
	OrderDetails OrderDetailList `json:"order_details" binding:"required"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

func (p *Order) CreateUUID() string {
	newUUID, _ := uuid.NewRandom()
	return newUUID.String()
}
