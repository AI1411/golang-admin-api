package handler

import (
	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type OrderHandler struct {
	Db *gorm.DB
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{Db: db}
}

type searchOrderPrams struct {
	ID         *string `form:"id" binding:"omitempty"`
	UserId     *string `form:"user_id" binding:"omitempty"`
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
	order.CreateUUID()
	order.OrderStatus = models.OrderStatusNew
	if err := h.Db.Create(&order).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to create order", err))
		return
	}
	ctx.JSON(http.StatusCreated, order)
	return
}
