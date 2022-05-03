package models

import "github.com/google/uuid"

type OrderDetailList []OrderDetail

type OrderDetail struct {
	ID        string `json:"id"`
	OrderID   string `json:"order_id" binding:"required,len=36"`
	ProductID string `json:"product_id" binding:"required,len=36"`
	Quantity  int64  `json:"quantity" binding:"required,gte=1"`
	Price     int64  `json:"price" binding:"required,gte=1"`
}

func (l *OrderDetailList) TotalPrice() int64 {
	var totalPrice int64
	for _, v := range *l {
		totalPrice = v.Quantity * v.Price
	}
	return totalPrice
}

func (d *OrderDetail) CreateUUID() {
	newUUID, _ := uuid.NewRandom()
	d.ID = newUUID.String()
}
