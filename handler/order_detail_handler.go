package handler

import (
	"net/http"

	"github.com/AI1411/golang-admin-api/util/appcontext"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"
)

type OrderDetailHandler struct {
	Db     *gorm.DB
	logger *zap.Logger
}

func NewOrderDetailHandler(db *gorm.DB, logger *zap.Logger) *OrderDetailHandler {
	return &OrderDetailHandler{
		Db:     db,
		logger: logger,
	}
}

type orderDetailRequest struct {
	OrderId   string `json:"order_id" binding:"required,uuid4" format:"/^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i" example:"e6f8b0c0-c8b7-4c8e-a098-e1c8b9e1c9c9"`
	ProductId string `json:"product_id" binding:"required,uuid4" format:"/^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i" example:"e6f8b0c0-c8b7-4c8e-a098-e1c8b9e1c9c9"`
	Quantity  int    `json:"quantity" binding:"required,min=1" example:"1"`
	Price     int    `json:"price" binding:"required,min=1" example:"1"`
}

type orderDetailResponseItem struct {
	Id                string `json:"id" example:"e6f8b0c0-c8b7-4c8e-a098-e1c8b9e1c9c9"`
	OrderId           string `json:"order_id" example:"e6f8b0c0-c8b7-4c8e-a098-e1c8b9e1c9c9"`
	ProductId         string `json:"product_id" example:"e6f8b0c0-c8b7-4c8e-a098-e1c8b9e1c9c9"`
	Quantity          int    `json:"quantity" example:"1"`
	OrderDetailStatus string `json:"order_detail_status" example:"new"`
	Price             int    `json:"price" example:"1"`
}

// GetOrderDetail @title 注文明細詳細取得
// @id GetOrderDetail
// @tags orderDetails
// @version バージョン(1.0)
// @description 指定された条件に一致するorderDetail詳細情報を取得する
// @Summary orderDetail詳細取得
// @Produce json
// @Success 200 {object} orderDetailResponseItem
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /orderDetails [GET]
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
func (h *OrderDetailHandler) GetOrderDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	traceID := appcontext.GetTraceID(ctx)
	var orderDetail models.OrderDetail
	if err := h.Db.Where("id = ?", id).First(&orderDetail).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			h.logger.Error("failed to get order detail", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("order detail not found"))
		case gorm.ErrInvalidSQL:
			h.logger.Error("invalid sql", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	ctx.JSON(http.StatusOK, orderDetail)
}

// CreateOrderDetail @title orderDetail作成
// @id CreateOrderDetail
// @tags orderDetails
// @version バージョン(1.0)
// @description orderDetailを作成する
// @Summary orderDetail作成
// @Produce json
// @Success 201 {object} orderDetailResponseItem
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /orderDetails [POST]
// @Accept json
// @Param orderDetailRequest body orderDetailRequest true "create orderDetail"
func (h *OrderDetailHandler) CreateOrderDetail(ctx *gin.Context) {
	traceID := appcontext.GetTraceID(ctx)
	orderDetail := models.OrderDetail{}
	if err := ctx.ShouldBindJSON(&orderDetail); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	orderDetail.CreateUUID()
	if err := h.Db.Create(&orderDetail).Error; err != nil {
		h.logger.Error("failed to create order detail", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.JSON(http.StatusInternalServerError,
			errors.NewInternalServerError("failed to create order detail", err))
		return
	}
	ctx.JSON(http.StatusCreated, orderDetail)
}

// UpdateOrderDetail @title orderDetail編集
// @id UpdateOrderDetail
// @tags orderDetails
// @version バージョン(1.0)
// @description orderDetailを編集する
// @Summary orderDetail編集
// @Produce json
// @Success 202 {object} orderDetailResponseItem
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /orderDetails/:id [PUT]
// @Accept json
// @Param orderDetailRequest body orderDetailRequest true "update orderDetail"
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
func (h *OrderDetailHandler) UpdateOrderDetail(ctx *gin.Context) {
	var orderDetail models.OrderDetail
	id := ctx.Param("id")
	traceID := appcontext.GetTraceID(ctx)
	if err := h.Db.Where("id = ?", id).First(&orderDetail).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			h.logger.Error("failed to update order detail", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("order detail not found"))
		case gorm.ErrInvalidSQL:
			h.logger.Error("invalid sql", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	if err := ctx.ShouldBindJSON(&orderDetail); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	if err := h.Db.Save(&orderDetail).Error; err != nil {
		h.logger.Error("failed to update order detail", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.JSON(http.StatusInternalServerError,
			errors.NewInternalServerError("failed to update order detail", err))
		return
	}

	ctx.JSON(http.StatusAccepted, orderDetail)
}

// DeleteOrderDetail @title orderDetail削除
// @id DeleteOrderDetail
// @tags orderDetails
// @version バージョン(1.0)
// @description orderDetailを削除する
// @Summary orderDetail削除
// @Produce json
// @Success 204
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /orderDetails/:id [DELETE]
// @Accept json
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
func (h *OrderDetailHandler) DeleteOrderDetail(ctx *gin.Context) {
	orderDetail := models.OrderDetail{}
	id := ctx.Param("id")
	traceID := appcontext.GetTraceID(ctx)
	if err := h.Db.Where("id = ?", id).First(&orderDetail).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			h.logger.Error("failed to delete order detail", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("order detail not found"))
		case gorm.ErrInvalidSQL:
			h.logger.Error("invalid sql", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	if err := h.Db.Delete(&orderDetail).Error; err != nil {
		h.logger.Error("failed to delete order detail", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.JSON(http.StatusInternalServerError,
			errors.NewInternalServerError("failed to delete order detail", err))
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
