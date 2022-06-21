package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"
)

type ProductHandler struct {
	Db *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{Db: db}
}

type searchProductParams struct {
	ProductName string `form:"product_name" binding:"max=64"`
	PriceFrom   string `form:"price_from" binding:"omitempty,numeric"`
	PriceTo     string `form:"price_to" binding:"omitempty,numeric"`
	Remarks     string `form:"remarks" binding:"omitempty,max=255"`
	Quantity    string `form:"quantity" binding:"omitempty,numeric,min=1"`
	Offset      string `form:"offset,default=0" binding:"omitempty,numeric"`
	Limit       string `form:"limit,default=10" binding:"omitempty,numeric"`
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
	if err := query.Find(&products).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to get products", err))
		return
	}
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
}

func (h *ProductHandler) UpdateProduct(ctx *gin.Context) {
	product := models.Product{}
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
	if err := ctx.ShouldBindJSON(&product); err != nil {
		restErr := errors.NewBadRequestError("invalid request")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	if err := h.Db.Save(&product).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to update product", err))
		return
	}
	ctx.JSON(http.StatusAccepted, product)
}

func (h *ProductHandler) DeleteProduct(ctx *gin.Context) {
	product := models.Product{}
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
	if err := h.Db.Delete(&product).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to delete product", err))
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func createProductQueryBuilder(params searchProductParams, h *ProductHandler) *gorm.DB {
	var products []models.Product
	query := h.Db.Find(&products)
	if params.ProductName != "" {
		query = query.Where("product_name LIKE ?", "%"+params.ProductName+"%")
	}
	if params.PriceFrom != "" {
		query = query.Where("price > ?", params.PriceFrom)
	}
	if params.PriceTo != "" {
		query = query.Where("price < ?", params.PriceTo)
	}
	if params.Remarks != "" {
		query = query.Where("remarks LIKE ?", "%"+params.Remarks+"%")
	}
	if params.Quantity != "" {
		query = query.Where("quantity = ?", params.Quantity)
	}
	if params.Offset != "" {
		query = query.Offset(params.Offset)
	}
	if params.Limit != "" {
		query = query.Limit(params.Limit)
	}
	return query
}
