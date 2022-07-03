package models

type CouponUser struct {
	ID       int64  `json:"id"`
	CouponID string `json:"coupon_id" binding:"required"`
	UserID   string `json:"user_id" binding:"required"`
	UseCount uint   `json:"use_count"`
}
