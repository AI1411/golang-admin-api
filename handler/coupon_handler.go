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

	if err := h.Db.Table("coupon_user").Where("coupon_id = ? and user_id = ?", couponID, userID).First(&couponUser).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to find coupon_user", err))
		return
	}

	if couponUser.UseCount == 0 {
		ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("coupon already acquired"))
		return
	}

	if err := h.Db.Table("coupon_user").Create(&couponUser).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to acquire coupon", err))
		return
	}

	ctx.JSON(http.StatusCreated, couponUser)
}
