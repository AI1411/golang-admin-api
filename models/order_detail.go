package models

import "github.com/google/uuid"

type OrderDetailList []OrderDetail

type OrderDetailStatus string

const (
	OrderDetailStatusNew       OrderDetailStatus = "new"
	OrderDetailStatusPaid      OrderDetailStatus = "paid"
	OrderDetailStatusCanceled  OrderDetailStatus = "canceled"
	OrderDetailStatusDelivered OrderDetailStatus = "delivered"
	OrderDetailStatusRefunded  OrderDetailStatus = "refunded"
	OrderDetailStatusReturned  OrderDetailStatus = "returned"
)

type OrderDetail struct {
	ID                string            `json:"id"`
	OrderID           string            `json:"order_id" binding:"omitempty,len=36"`
	ProductID         string            `json:"product_id" binding:"required,len=36"`
	Quantity          int64             `json:"quantity" binding:"required,gte=1"`
	OrderDetailStatus OrderDetailStatus `json:"order_detail_status" binding:"omitempty,oneof=new paid cancelled delivered refunded returned"`
	Price             int64             `json:"price" binding:"required,gte=1"`
}

func (l *OrderDetailList) TotalPrice() int64 {
	var totalPrice int64
	for _, v := range *l {
		totalPrice += v.Quantity * v.Price
	}
	return totalPrice
}

func (l *OrderDetailList) TotalQuantity() int64 {
	var totalQuantity int64
	for _, v := range *l {
		totalQuantity += v.Quantity
	}
	return totalQuantity
}

func (d *OrderDetail) CreateUUID() {
	newUUID, _ := uuid.NewRandom()
	d.ID = newUUID.String()
}
