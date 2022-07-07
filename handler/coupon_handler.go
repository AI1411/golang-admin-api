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

type CouponHandler struct {
	Db     *gorm.DB
	logger *zap.Logger
}

func NewCouponHandler(db *gorm.DB, logger *zap.Logger) *CouponHandler {
	return &CouponHandler{
		Db:     db,
		logger: logger,
	}
}

type couponRequest struct {
	Title             string    `json:"title" binding:"required" example:"クーポンタイトル"`
	Remarks           string    `json:"remarks" binding:"omitempty,max=255" example:"クーポン備考"`
	DiscountAmount    uint64    `json:"discount_amount" binding:"omitempty,min=1" example:"1000"`
	DiscountRate      uint8     `json:"discount_rate" binding:"omitempty,min=1,max=100" example:"10"`
	MaxDiscountAmount uint64    `json:"max_discount_amount" binding:"omitempty,min=1" example:"1000"`
	UseStartAt        time.Time `json:"use_start_at" binding:"required" example:"2020-01-01T00:00:00+09:00"`
	UseEndAt          time.Time `json:"use_end_at" binding:"required" example:"2020-01-01T00:00:00+09:00"`
	PublicStartAt     time.Time `json:"public_start_at" binding:"required" example:"2020-01-01T00:00:00+09:00"`
	PublicEndAt       time.Time `json:"public_end_at" binding:"required" example:"2020-01-01T00:00:00+09:00"`
	IsPublic          bool      `json:"is_public" example:"true"`
	IsPremium         bool      `json:"is_premium" example:"true"`
}

type acquireCouponRequest struct {
	UserID   string `json:"user_id" binding:"required,uuid4" example:"443b5f1c-8a3a-4485-b3bc-05e69b40b290" format:"/^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i"`
	CouponID string `json:"coupon_id" binding:"required,uuid4" example:"443b5f1c-8a3a-4485-b3bc-05e69b40b290" format:"/^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i"`
}

type couponDiscountRequest struct {
	CouponID   string   `json:"coupon_id" binding:"required,uuid4" example:"443b5f1c-8a3a-4485-b3bc-05e69b40b290" format:"/^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i"`
	UserID     string   `json:"user_id" binding:"required,uuid4" example:"443b5f1c-8a3a-4485-b3bc-05e69b40b290" format:"/^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i"`
	ProductIDs []string `json:"product_ids" binding:"required" example:"443b5f1c-8a3a-4485-b3bc-05e69b40b290" format:"/^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i"`
}

type couponResponseItem struct {
	models.Coupon
}

type couponResponse struct {
	Coupons []couponResponseItem `json:"coupons"`
	Total   int                  `json:"total"`
}

type discountedProductResponseItem struct {
	ProductName string `json:"product_name" example:"商品名" format:"/^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i"`
	Price       uint   `json:"price" example:"1000"`
}

type searchCouponParams struct {
	Title             string `form:"title" binding:"omitempty,max=64"`
	DiscountAmount    string `form:"discount_amount" binding:"omitempty,numeric"`
	DiscountRate      string `form:"discount_rate" binding:"omitempty,numeric"`
	MaxDiscountAmount string `form:"max_discount_amount" binding:"omitempty,numeric"`
	UseStartAtFrom    string `form:"use_start_at_from" binding:"omitempty,datetime"`
	UseStartAtTo      string `form:"use_start_at_to" binding:"omitempty,datetime"`
	UseEndAtFrom      string `form:"use_end_at_from" binding:"omitempty,datetime"`
	UseEndAtTo        string `form:"use_end_at_to" binding:"omitempty,datetime"`
	PublicStartAtFrom string `form:"public_start_at_from" binding:"omitempty,datetime"`
	PublicStartAtTo   string `form:"public_start_at_to" binding:"omitempty,datetime"`
	PublicEndAtFrom   string `form:"public_end_at_from" binding:"omitempty,datetime"`
	PublicEndAtTo     string `form:"public_end_at_to" binding:"omitempty,datetime"`
	IsPublic          string `form:"is_public" binding:"omitempty,boolean"`
	IsPremium         string `form:"is_premium" binding:"omitempty,boolean"`
	Offset            string `form:"offset,default=0" binding:"omitempty,numeric"`
	Limit             string `form:"limit,default=10" binding:"omitempty,numeric"`
}

// GetAllCoupon @title 一覧取得
// @id GetAllCoupon
// @tags coupons
// @version バージョン(1.0)
// @description 指定された条件に一致するcoupon一覧情報を取得する
// @Summary coupon一覧取得
// @Produce json
// @Success 200 {object} couponResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /coupons [GET]
// @Param title query string false "タイトル" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param discount_amount query int false "値引額" minimum(0)
// @Param discount_rate query int false "割引率" minimum(0)
// @Param max_discount_amount query int false "最大値引額" minimum(0)
// @Param use_start_at_from query string false "利用開始日" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param use_start_at_to query string false "利用開始日" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param use_end_at_from query string false "利用終了日" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param use_end_at_to query string false "利用終了日" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param public_start_at_from query string false "公開開始日" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param public_start_at_to query string false "公開開始日" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param public_end_at_from query string false "公開終了日" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param public_end_at_to query string false "公開終了日" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param is_public query string false "公開フラグ"
// @Param is_premium query string false "プレミアムフラグ"
// @Param offset query int false "開始位置" default(0) minimum(0)
// @Param limit query int false "取得上限" default(12) minimum(1) maximum(100)
func (h *CouponHandler) GetAllCoupon(ctx *gin.Context) {
	traceID := appcontext.GetTraceID(ctx)
	var params searchCouponParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind query params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	var coupons []models.Coupon
	query := createCouponQueryBuilder(params, h)
	if err := query.Find(&coupons).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to get coupons", err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"total":   len(coupons),
		"coupons": coupons,
	})
}

// GetCouponDetail @title coupon詳細
// @id GetCouponDetail
// @tags coupons
// @version バージョン(1.0)
// @description coupon詳細を返す
// @Summary coupon詳細取得
// @Produce json
// @Success 200 {object} couponResponseItem
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /coupons/:id [GET]
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
func (h *CouponHandler) GetCouponDetail(ctx *gin.Context) {
	var coupon models.Coupon
	id := ctx.Param("id")
	traceID := appcontext.GetTraceID(ctx)
	if err := h.Db.Where("id = ?", id).First(&coupon).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			h.logger.Error("failed to find coupon", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("coupon not found"))
		case gorm.ErrInvalidSQL:
			h.logger.Error("invalid sql", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	ctx.JSON(http.StatusOK, coupon)
}

// CreateCoupon @title coupon作成
// @id CreateCoupon
// @tags coupons
// @version バージョン(1.0)
// @description couponを作成する
// @Summary coupon作成
// @Produce json
// @Success 201 {object} couponResponseItem
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /coupons [POST]
// @Accept json
// @Param couponRequest body couponRequest true "create coupon"
func (h *CouponHandler) CreateCoupon(ctx *gin.Context) {
	coupon := models.Coupon{}
	if err := ctx.ShouldBindJSON(&coupon); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	traceID := appcontext.GetTraceID(ctx)
	coupon.CreateUUID()
	if err := h.Db.Create(&coupon).Error; err != nil {
		h.logger.Error("failed to create coupon", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to create coupon", err))
		return
	}
	ctx.JSON(http.StatusCreated, coupon)
}

// UpdateCoupon @title coupon編集
// @id UpdateCoupon
// @tags coupons
// @version バージョン(1.0)
// @description couponを編集する
// @Summary coupon編集
// @Produce json
// @Success 202 {object} couponResponseItem
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /coupons/:id [PUT]
// @Accept json
// @Param couponRequest body couponRequest true "update coupon"
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
func (h *CouponHandler) UpdateCoupon(ctx *gin.Context) {
	coupon := models.Coupon{}
	id := ctx.Param("id")
	traceID := appcontext.GetTraceID(ctx)
	if err := h.Db.Where("id = ?", id).First(&coupon).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			h.logger.Error("failed to update coupon", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("coupon not found"))
		case gorm.ErrInvalidSQL:
			h.logger.Error("invalid sql", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	if err := ctx.ShouldBindJSON(&coupon); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	if err := h.Db.Save(&coupon).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to update coupon", err))
		return
	}
	ctx.JSON(http.StatusAccepted, coupon)
}

// AcquireCoupon @title coupon獲得
// @id AcquireCoupon
// @tags coupons
// @version バージョン(1.0)
// @description couponを獲得する
// @Summary coupon獲得
// @Produce json
// @Success 201
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /:coupon_id/users/:user_id [POST]
// @Accept json
// @Param acquireCouponRequest body acquireCouponRequest true "acquire coupon"
func (h *CouponHandler) AcquireCoupon(ctx *gin.Context) {
	var req acquireCouponRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	traceID := appcontext.GetTraceID(ctx)
	if err := h.Db.Transaction(func(tx *gorm.DB) error {
		var coupon models.Coupon
		if err := h.Db.First(&coupon, "id = ?", req.CouponID).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				h.logger.Error("failed to acquire coupon", zap.Error(err),
					zap.String("trace_id", traceID))
				ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("coupon not found"))
				return err
			}
		}

		var user models.User
		if err := h.Db.First(&user, "id = ?", req.UserID).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				h.logger.Error("failed to find user", zap.Error(err),
					zap.String("trace_id", traceID))
				ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("user not found"))
				return err
			}
		}

		var couponUser models.CouponUser
		if err := h.Db.Table("coupon_user").
			First(&couponUser, "coupon_id = ? and user_id = ?", req.CouponID, req.UserID).
			Error; err != nil {
			h.logger.Error("failed to find coupon user", zap.Error(err),
				zap.String("trace_id", traceID))
			return err
		}

		if couponUser.ID != 0 {
			h.logger.Error("coupon already acquired", zap.Error(nil),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusBadRequest,
				errors.NewBadRequestError("failed to acquire coupon"))
			return nil
		}

		if err := h.Db.Table("coupon_user").Create(&models.CouponUser{
			CouponID: req.CouponID,
			UserID:   req.UserID,
			UseCount: 0,
		}).Error; err != nil {
			h.logger.Error("failed to acquire coupon", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusInternalServerError,
				errors.NewInternalServerError("failed to acquire coupon", err))
			return err
		}
		return nil
	}); err != nil {
		h.logger.Error("failed to acquire coupon", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to acquire coupon", err))
		return
	}

	ctx.Status(http.StatusCreated)
}

func (h *CouponHandler) DiscountedList(ctx *gin.Context) {
	var req couponDiscountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	traceID := appcontext.GetTraceID(ctx)
	var coupon models.Coupon
	if err := h.Db.First(&coupon, "id = ?", req.CouponID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			h.logger.Error("failed to get coupon", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("coupon not found"))
			return
		}
	}

	var user models.User
	if err := h.Db.First(&user, "id = ?", req.UserID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			h.logger.Error("failed to find user", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("user not found"))
			return
		}
	}

	var couponUser models.CouponUser
	if err := h.Db.Table("coupon_user").
		First(&couponUser, "coupon_id = ? and user_id = ?", req.CouponID, req.UserID).
		Error; err != nil {
		h.logger.Error("failed to find coupon user", zap.Error(err),
			zap.String("trace_id", traceID))
		return
	}

	var products []models.Product
	if err := h.Db.Where(req.ProductIDs).Find(&products).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			h.logger.Error("failed to find products", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("products not found"))
			return
		}
	}

	if len(products) == 0 {
		h.logger.Error("no products", zap.Error(nil), zap.String("trace_id", traceID))
		ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("products not found"))
		return
	}

	for i, product := range products {
		products[i].Price = product.Price - uint(coupon.DiscountAmount)
	}
	ctx.JSON(http.StatusOK, products)
	return
}

func createCouponQueryBuilder(params searchCouponParams, h *CouponHandler) *gorm.DB {
	var coupons []models.Coupon
	query := h.Db.Find(&coupons)

	if params.Title != "" {
		query = query.Where("title LIKE ?", "%"+params.Title+"%")
	}
	if params.DiscountAmount != "" {
		query = query.Where("discount_amount = ?", params.DiscountAmount)
	}
	if params.DiscountRate != "" {
		query = query.Where("discount_rate = ?", params.DiscountRate)
	}
	if params.MaxDiscountAmount != "" {
		query = query.Where("max_discount_amount = ?", params.MaxDiscountAmount)
	}
	if params.UseStartAtFrom != "" {
		query = query.Where("use_start_at > ?", params.UseStartAtFrom)
	}
	if params.UseStartAtTo != "" {
		query = query.Where("use_start_at < ?", params.UseStartAtTo)
	}
	if params.UseEndAtFrom != "" {
		query = query.Where("use_end_at > ?", params.UseEndAtFrom)
	}
	if params.UseEndAtTo != "" {
		query = query.Where("use_end_at < ?", params.UseEndAtTo)
	}
	if params.PublicStartAtFrom != "" {
		query = query.Where("public_start_at > ?", params.PublicStartAtFrom)
	}
	if params.PublicStartAtTo != "" {
		query = query.Where("public_start_at < ?", params.PublicStartAtTo)
	}
	if params.PublicEndAtFrom != "" {
		query = query.Where("public_end_at > ?", params.PublicEndAtFrom)
	}
	if params.PublicEndAtTo != "" {
		query = query.Where("public_end_at < ?", params.PublicEndAtTo)
	}
	if params.IsPublic != "" {
		query = query.Where("is_public = ?", params.IsPublic)
	}
	if params.IsPremium != "" {
		query = query.Where("is_premium = ?", params.IsPremium)
	}
	if params.Offset != "" {
		query = query.Offset(params.Offset)
	}
	if params.Limit != "" {
		query = query.Limit(params.Limit)
	}
	return query
}
