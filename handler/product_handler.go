package handler

import (
	"github.com/AI1411/golang-admin-api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type ProductHandler struct {
	Db *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{Db: db}
}

type searchProductParams struct {
	ProductID   *uuid.UUID `json:"product_id"`
	ProductName *string    `json:"product_name"`
	Price       *string    `json:"price"`
	Remarks     *string    `json:"remarks"`
	Quantity    *string    `json:"quantity"`
	Offset      string     `json:"offset"`
	Limit       string     `json:"limit"`
}

func (h ProductHandler) GetAllProduct(ctx *gin.Context) {
	var params searchProductParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	var products []models.Product
	h.Db.Find(&products)
	ctx.JSON(200, products)
}

func createProductQueryBuilder(params searchProductParams, h *ProductHandler) *gorm.DB {
	var products []models.Product
	query := h.Db.Find(&products)

	if params.ProductID != nil {
		query = query.Where("product_id = ?", params.ProductID)
	}
	if params.ProductName != nil {
		query = query.Where("product_name LIKE ?", "%"+*params.ProductName+"%")
	}
	if params.Price != nil {
		query = query.Where("price = ?", *params.Price)
	}
	if params.Remarks != nil {
		query = query.Where("remarks LIKE ?", "%"+*params.Remarks+"%")
	}
	if params.Quantity != nil {
		query = query.Where("quantity = ?", *params.Quantity)
	}
	if params.Offset != "" {
		query = query.Offset(params.Offset)
	}
	if params.Limit != "" {
		query = query.Limit(params.Limit)
	}
	return query
}
