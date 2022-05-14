package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"
)

type OrderDetailHandler struct {
	Db *gorm.DB
}

func NewOrderDetailHandler(db *gorm.DB) *OrderDetailHandler {
	return &OrderDetailHandler{Db: db}
}

func (h *OrderDetailHandler) CreateOrderDetail(ctx *gin.Context) {
	orderDetail := models.OrderDetail{}
	if err := ctx.ShouldBindJSON(&orderDetail); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	orderDetail.CreateUUID()
	if err := h.Db.Create(&orderDetail).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to create order detail", err))
		return
	}
	ctx.JSON(http.StatusCreated, orderDetail)
}
