package handler

import (
	"github.com/AI1411/golang-admin-api/util/appcontext"
	"go.uber.org/zap"
	"net/http"

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

func (h *CouponHandler) GetCouponDetail(ctx *gin.Context) {
	var coupon models.Coupon
	id := ctx.Param("id")
	if err := h.Db.Where("id = ?", id).First(&coupon).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("coupon not found"))
		case gorm.ErrInvalidSQL:
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	ctx.JSON(http.StatusOK, coupon)
}

func (h *CouponHandler) CreateCoupon(ctx *gin.Context) {
	coupon := models.Coupon{}
	if err := ctx.ShouldBindJSON(&coupon); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	coupon.CreateUUID()
	if err := h.Db.Create(&coupon).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to create coupon", err))
		return
	}
	ctx.JSON(http.StatusCreated, coupon)
}

func (h *CouponHandler) UpdateCoupon(ctx *gin.Context) {
	coupon := models.Coupon{}
	id := ctx.Param("id")
	if err := h.Db.Where("id = ?", id).First(&coupon).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("coupon not found"))
		case gorm.ErrInvalidSQL:
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

func (h *CouponHandler) AcquireCoupon(ctx *gin.Context) {
	var couponUser models.CouponUser
	couponID := ctx.Param("coupon_id")
	userID := ctx.Param("user_id")
	couponUser.CouponID = couponID
	couponUser.UserID = userID
	if err := ctx.ShouldBindJSON(&couponUser); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	if err := h.Db.Transaction(func(tx *gorm.DB) error {
		var coupon models.Coupon
		if err := h.Db.First(&coupon, "id = ?", couponID).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("coupon not found"))
				return err
			}
		}

		var user models.User
		if err := h.Db.First(&user, "id = ?", userID).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("user not found"))
				return err
			}
		}

		if err := h.Db.Table("coupon_user").
			Where("coupon_id = ? and user_id = ?", couponID, userID).
			First(&couponUser).Error; err != nil {
			return err
		}

		if err := h.Db.Table("coupon_user").Create(&couponUser).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError,
				errors.NewInternalServerError("failed to acquire coupon", err))
			return err
		}
		return nil
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to acquire coupon", err))
		return
	}

	ctx.JSON(http.StatusCreated, couponUser)
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
