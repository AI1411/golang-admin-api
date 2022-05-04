package handler

import (
	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type OrderHandler struct {
	Db *gorm.DB
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{Db: db}
}

type searchOrderPrams struct {
	ID         *string `form:"id" binding:"omitempty"`
	UserID     *string `form:"user_id" binding:"omitempty"`
	Quantity   *string `form:"quantity" binding:"omitempty"`
	TotalPrice *string `form:"total_price" binding:"omitempty,numeric"`
	CreatedAt  *string `form:"created_at" binding:"omitempty,datetime"`
	Offset     string  `form:"offset" binding:"omitempty,numeric"`
	Limit      string  `form:"limit" binding:"omitempty,numeric"`
}

func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	order := models.Order{}
	if err := ctx.ShouldBindJSON(&order); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	orderData := models.Order{
		ID:          order.CreateUUID(),
		UserID:      order.UserID,
		Quantity:    order.OrderDetails.TotalQuantity(),
		TotalPrice:  order.OrderDetails.TotalPrice(),
		OrderStatus: models.OrderStatusNew,
		Remarks:     order.Remarks,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := h.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("orders").Create(&orderData).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to create order", err))
			return err
		}
		for i, detail := range order.OrderDetails {
			detail.ID = order.CreateUUID()
			order.OrderDetails[i].OrderID = orderData.ID
			order.OrderDetails[i].OrderID = orderData.ID
			orderData.OrderDetails = append(orderData.OrderDetails, detail)
			if err := tx.Create(&detail).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to create order detail", err))
				return err
			}
		}
		return nil
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to create order", err))
		return
	}

	ctx.JSON(http.StatusCreated, orderData)
	return
}
