package models

import (
	"github.com/google/uuid"
	"time"
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
	ID          uuid.UUID   `json:"id"`
	UserId      int64       `json:"user_id" binding:"required,min=1"`
	Quantity    int64       `json:"quantity" binding:"required,gte=1"`
	TotalPrice  int64       `json:"total_price" binding:"required,min=0"`
	OrderStatus OrderStatus `json:"order_status" binding:"omitempty,oneof=new paid cancelled delivered refunded returned partially partially_paid"`
	Remarks     string      `json:"remarks" binding:"omitempty,max=255"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func (p *Order) CreateUUID() {
	p.ID = uuid.New()
}
