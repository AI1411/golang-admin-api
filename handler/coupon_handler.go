package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"
)

type CouponHandler struct {
	Db *gorm.DB
}

func NewCouponHandler(db *gorm.DB) *CouponHandler {
	return &CouponHandler{Db: db}
}

type searchCouponParams struct {
	ID                string `json:"id"`
	Title             string `json:"title"`
	Remarks           string `json:"remarks"`
	DiscountAmount    string `json:"discountAmount"`
	DiscountRate      string `json:"discountRate"`
	MaxDiscountAmount string `json:"maxDiscountAmount"`
	UseStartAtFrom    string `json:"useStartAtFrom"`
	UseStartAtTo      string `json:"useStartAtTo"`
	UseEndAtFrom      string `json:"useEndAtFrom"`
	UseEndAtTo        string `json:"useEndAtTo"`
	PublicStartAtFrom string `json:"publicStartAtFrom"`
	PublicStartAtTo   string `json:"publicStartAtTo"`
	PublicEndAtFrom   string `json:"publicEndAtFrom"`
	PublicEndAtTo     string `json:"publicEndAtTo"`
	IsPublic          string `json:"isPublic"`
	IsPremium         string `json:"isPremium"`
	Offset            string `json:"offset"`
	Limit             string `json:"limit"`
}

func (h *CouponHandler) GetAllCoupon(ctx *gin.Context) {
	var params searchCouponParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		res := createValidateErrorResponse(err)
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

	if params.ID != "" {
		query = query.Where("id = ?", params.ID)
	}
	if params.Title != "" {
		query = query.Where("title LIKE ?", "%"+params.Title+"%")
	}
	if params.Remarks != "" {
		query = query.Where("title LIKE ?", "%"+params.Remarks+"%")
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
		query = query.Where("quantity > ?", params.UseStartAtFrom)
	}
	if params.UseStartAtTo != "" {
		query = query.Where("quantity < ?", params.UseStartAtTo)
	}
	if params.UseEndAtFrom != "" {
		query = query.Where("quantity > ?", params.UseEndAtFrom)
	}
	if params.UseEndAtTo != "" {
		query = query.Where("quantity < ?", params.UseEndAtTo)
	}
	if params.PublicStartAtFrom != "" {
		query = query.Where("quantity > ?", params.PublicStartAtFrom)
	}
	if params.PublicStartAtTo != "" {
		query = query.Where("quantity < ?", params.PublicStartAtTo)
	}
	if params.PublicEndAtFrom != "" {
		query = query.Where("quantity > ?", params.PublicEndAtFrom)
	}
	if params.PublicEndAtTo != "" {
		query = query.Where("quantity < ?", params.PublicEndAtTo)
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
