package models

type CouponUser struct {
	ID       int64  `json:"id"`
	CouponID string `json:"coupon_id"`
	UserID   string `json:"user_id"`
	UseCount uint   `json:"use_count"`
}
