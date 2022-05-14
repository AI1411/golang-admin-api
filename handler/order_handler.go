package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"
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

func (h *OrderHandler) GetOrders(ctx *gin.Context) {
	var params searchOrderPrams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	var orders []models.Order
	query := createOrderQueryBuilder(params, h)
	query.Preload("OrderDetails").Find(&orders)

	ctx.JSON(http.StatusOK, gin.H{
		"total":  len(orders),
		"orders": orders,
	})
}

func (h *OrderHandler) GetOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	var order models.Order
	if err := h.Db.Preload("OrderDetails").Where("id = ?", id).First(&order).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, errors.NewNotFoundError("order not found"))
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to get order", err))
		return
	}

	ctx.JSON(http.StatusOK, order)
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
		for _, detail := range order.OrderDetails {
			detail.ID = order.CreateUUID()
			detail.OrderID = orderData.ID
			detail.OrderDetailStatus = models.OrderDetailStatusNew
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
}

func createOrderQueryBuilder(params searchOrderPrams, h *OrderHandler) *gorm.DB {
	var orders []models.Order
	query := h.Db.Find(&orders)
	if params.ID != nil {
		query = query.Where("id = ?", *params.ID)
	}
	if params.UserID != nil {
		query = query.Where("user_id = ?", *params.UserID)
	}
	if params.Quantity != nil {
		query = query.Where("quantity = ?", *params.Quantity)
	}
	if params.CreatedAt != nil {
		query = query.Where("created_at = ?", *params.CreatedAt)
	}
	if params.Offset != "" {
		query = query.Offset(params.Offset)
	}
	if params.Limit != "" {
		query = query.Limit(params.Limit)
	}
	return query
}
