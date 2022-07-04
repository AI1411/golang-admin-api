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

type ProductHandler struct {
	Db            *gorm.DB
	logger        *zap.Logger
	uuidGenerator models.UUIDGenerator
}

func NewProductHandler(db *gorm.DB, uuidGenerator models.UUIDGenerator, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{
		Db:            db,
		logger:        logger,
		uuidGenerator: uuidGenerator,
	}
}

type searchProductParams struct {
	ProductName string `form:"product_name" binding:"max=64"`
	PriceFrom   string `form:"price_from" binding:"omitempty,numeric"`
	PriceTo     string `form:"price_to" binding:"omitempty,numeric"`
	Quantity    string `form:"quantity" binding:"omitempty,numeric,min=1"`
	Offset      string `form:"offset,default=0" binding:"omitempty,numeric"`
	Limit       string `form:"limit,default=10" binding:"omitempty,numeric"`
}

type productRequest struct {
	ProductName string `json:"product_name" example:"product name" binding:"required,max=64"`
	Remarks     string `json:"remarks" example:"remarks" binding:"omitempty,max=255"`
}

type productResponseItem struct {
	Id          string `json:"id" example:"218c51c0-904e-4743-a2ae-94f0e34a0d6f"`
	ProductName string `json:"product_name" example:"product name"`
	Price       int    `json:"price" example:"100"`
	Remarks     string `json:"remarks" example:"remarks"`
	Quantity    int    `json:"quantity" example:"1"`
}

type productsResponse struct {
	Total    int                   `json:"total"`
	Products []productResponseItem `json:"products"`
}

// GetAllProduct @title 一覧取得
// @id GetAllProduct
// @tags products
// @version バージョン(1.0)
// @description 指定された条件に一致するproduct一覧情報を取得する
// @Summary product一覧取得
// @Produce json
// @Success 200 {object} productsResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /products [GET]
// @Param product_name query string false "商品名" maxlength(64)
// @Param price_from query int false "金額from" minimum(0)
// @Param price_to query int false "金額to" minimum(0)
// @Param quantity query int false "数量" minimum(1)
// @Param offset query int false "開始位置" default(0) minimum(0)
// @Param limit query int false "取得上限" default(12) minimum(1) maximum(100)
func (h *ProductHandler) GetAllProduct(ctx *gin.Context) {
	traceID := appcontext.GetTraceID(ctx)
	var params searchProductParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	var products []models.Product
	query := createProductQueryBuilder(params, h)
	if err := query.Find(&products).Error; err != nil {
		h.logger.Error("failed to get products", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to get products", err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"total":    len(products),
		"products": products,
	})
}

// GetProductDetail @title product詳細
// @id GetProductDetail
// @tags products
// @version バージョン(1.0)
// @description product詳細を返す
// @Summary product詳細取得
// @Produce json
// @Success 200 {object} productResponseItem
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /products/:id [GET]
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
func (h *ProductHandler) GetProductDetail(ctx *gin.Context) {
	var product models.Product
	id := ctx.Param("id")
	traceID := appcontext.GetTraceID(ctx)
	if err := h.Db.Where("id = ?", id).First(&product).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			h.logger.Error("failed to get product", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("product not found"))
		case gorm.ErrInvalidSQL:
			h.logger.Error("invalid sql", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	ctx.JSON(http.StatusOK, product)
}

// CreateProduct @title product作成
// @id CreateProduct
// @tags products
// @version バージョン(1.0)
// @description productを作成する
// @Summary product作成
// @Produce json
// @Success 201 {object} productResponseItem
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /products [POST]
// @Accept json
// @Param productRequest body productRequest true "create product"
func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	product := models.Product{}
	if err := ctx.ShouldBindJSON(&product); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	product.ID = h.uuidGenerator.GenerateUUID()
	traceID := appcontext.GetTraceID(ctx)
	if err := h.Db.Create(&product).Error; err != nil {
		h.logger.Error("failed to create product", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to create product", err))
		return
	}
	ctx.JSON(http.StatusCreated, product)
}

// UpdateProduct @title product編集
// @id UpdateProduct
// @tags products
// @version バージョン(1.0)
// @description productを編集する
// @Summary product編集
// @Produce json
// @Success 202 {object} productResponseItem
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /products/:id [PUT]
// @Accept json
// @Param productRequest body productRequest true "update product"
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
func (h *ProductHandler) UpdateProduct(ctx *gin.Context) {
	product := models.Product{}
	id := ctx.Param("id")
	traceID := appcontext.GetTraceID(ctx)
	if err := h.Db.Where("id = ?", id).First(&product).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			h.logger.Error("failed to update product", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("product not found"))
		case gorm.ErrInvalidSQL:
			h.logger.Error("invalid sql", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	if err := ctx.ShouldBindJSON(&product); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	if err := h.Db.Save(&product).Error; err != nil {
		h.logger.Error("failed to update product", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to update product", err))
		return
	}
	ctx.JSON(http.StatusAccepted, product)
}

// DeleteProduct @title product削除
// @id DeleteProduct
// @tags products
// @version バージョン(1.0)
// @description productを削除する
// @Summary product削除
// @Produce json
// @Success 204
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /products/:id [DELETE]
// @Accept json
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
func (h *ProductHandler) DeleteProduct(ctx *gin.Context) {
	product := models.Product{}
	id := ctx.Param("id")
	traceID := appcontext.GetTraceID(ctx)
	if err := h.Db.Where("id = ?", id).First(&product).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			h.logger.Error("failed to delete product", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("product not found"))
		case gorm.ErrInvalidSQL:
			h.logger.Error("invalid sql", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	if err := h.Db.Delete(&product).Error; err != nil {
		h.logger.Error("failed to delete product", zap.Error(err),
			zap.String("trace_id", traceID))
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
