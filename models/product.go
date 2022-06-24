package models

type Product struct {
	ID          string `json:"id"`
	ProductName string `json:"product_name" binding:"required,max=64"`
	Price       uint   `json:"price" binding:"required,gte=0"`
	Remarks     string `json:"remarks" binding:"omitempty,max=255"`
	Quantity    int    `json:"quantity" binding:"required,gte=0"`
}
