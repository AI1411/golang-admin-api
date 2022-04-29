package handler

import (
	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"net/http"
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

func (h *ProductHandler) GetAllProduct(ctx *gin.Context) {
	var params searchProductParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	var products []models.Product
	query := createProductQueryBuilder(params, h)
	query.Find(&products)
	ctx.JSON(http.StatusOK, gin.H{
		"total":    len(products),
		"products": products,
	})
}

func (h *ProductHandler) GetProductDetail(ctx *gin.Context) {
	var product models.Product
	id := ctx.Param("id")
	if err := h.Db.Where("id = ?", id).First(&product).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("product not found"))
		case gorm.ErrInvalidSQL:
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	ctx.JSON(http.StatusOK, product)
}

func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	product := models.Product{}
	if err := ctx.ShouldBindJSON(&product); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	product.CreateUUID()
	if err := h.Db.Create(&product).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to create product", err))
		return
	}
	ctx.JSON(http.StatusCreated, product)
	return
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
