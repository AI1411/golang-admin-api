package handler

import (
	"net/http"
	"time"

	"github.com/AI1411/golang-admin-api/util/appcontext"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"
)

type OrderHandler struct {
	Db     *gorm.DB
	logger *zap.Logger
}

func NewOrderHandler(db *gorm.DB, logger *zap.Logger) *OrderHandler {
	return &OrderHandler{
		Db:     db,
		logger: logger,
	}
}

type searchOrderParams struct {
	UserID      string `form:"user_id" binding:"omitempty,uuid4"`
	Quantity    string `form:"quantity" binding:"omitempty,numeric"`
	TotalPrice  string `form:"total_price" binding:"omitempty,numeric"`
	OrderStatus string `form:"order_status" binding:"omitempty,oneof=new paid canceled delivered refunded returned partially partially_paid"`
	Offset      string `form:"offset,default=0" binding:"omitempty,numeric"`
	Limit       string `form:"limit,default=10" binding:"omitempty,numeric"`
}

func (h *OrderHandler) GetOrders(ctx *gin.Context) {
	traceID := appcontext.GetTraceID(ctx)
	var params searchOrderParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
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
	traceID := appcontext.GetTraceID(ctx)
	order := models.Order{}
	if err := ctx.ShouldBindJSON(&order); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	orderData := models.Order{
		ID:          order.CreateUUID(),
		UserID:      order.UserID,
		Quantity:    order.OrderDetails.TotalQuantity(),
		TotalPrice:  order.OrderDetails.TotalPrice(),
		OrderStatus: order.OrderStatus,
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

func (h *OrderHandler) UpdateOrder(ctx *gin.Context) {
	order := models.Order{}
	id := ctx.Param("id")
	if err := h.Db.Where("id = ?", id).First(&order).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("order not found"))
		case gorm.ErrInvalidSQL:
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	traceID := appcontext.GetTraceID(ctx)
	if err := ctx.ShouldBindJSON(&order); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	if err := h.Db.Save(&order).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to update order", err))
		return
	}
	ctx.JSON(http.StatusAccepted, order)
}

func (h *OrderHandler) DeleteOrder(ctx *gin.Context) {
	var order models.Order
	id := ctx.Param("id")
	if err := h.Db.Where("id = ?", id).First(&order).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("order not found"))
		case gorm.ErrInvalidSQL:
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}

	if err := h.Db.Delete(&order).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to delete order", err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func createOrderQueryBuilder(params searchOrderParams, h *OrderHandler) *gorm.DB {
	var orders []models.Order
	query := h.Db.Order("created_at desc").Find(&orders)
	if params.UserID != "" {
		query = query.Where("user_id = ?", params.UserID)
	}
	if params.Quantity != "" {
		query = query.Where("quantity = ?", params.Quantity)
	}
	if params.TotalPrice != "" {
		query = query.Where("total_price = ?", params.TotalPrice)
	}
	if params.OrderStatus != "" {
		query = query.Where("order_status = ?", params.OrderStatus)
	}
	if params.Offset != "" {
		query = query.Offset(params.Offset)
	}
	if params.Limit != "" {
		query = query.Limit(params.Limit)
	}
	return query
}
