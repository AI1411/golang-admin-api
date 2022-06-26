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

func (h *OrderDetailHandler) GetOrderDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	var orderDetail models.OrderDetail
	if err := h.Db.Where("id = ?", id).First(&orderDetail).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("order detail not found"))
		case gorm.ErrInvalidSQL:
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	ctx.JSON(http.StatusOK, orderDetail)
}

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
		ctx.JSON(http.StatusInternalServerError,
			errors.NewInternalServerError("failed to create order detail", err))
		return
	}
	ctx.JSON(http.StatusCreated, orderDetail)
}

func (h *OrderDetailHandler) UpdateOrderDetail(ctx *gin.Context) {
	var orderDetail models.OrderDetail
	id := ctx.Param("id")
	if err := h.Db.Where("id = ?", id).First(&orderDetail).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("order detail not found"))
		case gorm.ErrInvalidSQL:
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	traceID := appcontext.GetTraceID(ctx)
	if err := ctx.ShouldBindJSON(&orderDetail); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	if err := h.Db.Save(&orderDetail).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError,
			errors.NewInternalServerError("failed to update order detail", err))
		return
	}

	ctx.JSON(http.StatusAccepted, orderDetail)
}

func (h *OrderDetailHandler) DeleteOrderDetail(ctx *gin.Context) {
	orderDetail := models.OrderDetail{}
	id := ctx.Param("id")
	if err := h.Db.Where("id = ?", id).First(&orderDetail).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("order detail not found"))
		case gorm.ErrInvalidSQL:
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	if err := h.Db.Delete(&orderDetail).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError,
			errors.NewInternalServerError("failed to delete order detail", err))
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
