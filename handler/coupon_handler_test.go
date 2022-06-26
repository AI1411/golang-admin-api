package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/AI1411/golang-admin-api/middleware"
	logger "github.com/AI1411/golang-admin-api/server"
	"github.com/gin-gonic/gin/binding"

	"github.com/stretchr/testify/require"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/AI1411/golang-admin-api/db"
)

const couponIDForTest = "090e142d-baa3-4039-9d21-cf5a1af39094"

var getCouponsTestCases = []struct {
	tid        int
	name       string
	request    map[string]interface{}
	wantStatus int
	wantBody   string
}{
	{
		tid:        1,
		name:       "クーポン一覧が正常に取得できること",
		request:    map[string]interface{}{},
		wantStatus: http.StatusOK,
		wantBody: `{
			"coupons": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"title": "coupon",
					"remarks": "coupon",
					"discount_amount": 1000,
					"discount_rate": 0,
					"max_discount_amount": 0,
					"use_start_at": "2022-06-01T00:00:00+09:00",
					"use_end_at": "2030-07-02T10:00:00+09:00",
					"public_start_at": "2022-06-01T00:00:00+09:00",
					"public_end_at": "2030-07-02T10:00:00+09:00",
					"is_public": false,
					"is_premium": false,
					"created_at": "2022-06-14T08:19:41+09:00",
					"updated_at": "2022-06-15T10:31:50+09:00"
				},
				{
					"id": "5c3325c1-d539-42d6-b405-2af2f6b99ed9",
					"title": "test",
					"remarks": "test",
					"discount_amount": 0,
					"discount_rate": 20,
					"max_discount_amount": 2000,
					"use_start_at": "2022-07-01T00:00:00+09:00",
					"use_end_at": "2032-07-02T10:00:00+09:00",
					"public_start_at": "2022-07-01T00:00:00+09:00",
					"public_end_at": "2032-07-02T10:00:00+09:00",
					"is_public": true,
					"is_premium": true,
					"created_at": "2022-06-14T08:19:12+09:00",
					"updated_at": "2022-06-14T08:19:12+09:00"
				}
			],
			"total": 2
		}`,
	},
	{
		tid:  2,
		name: "検索結果0件",
		request: map[string]interface{}{
			"title": "failed",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"coupons": [],
    		"total": 0
		}`,
	},
	{
		tid:  3,
		name: "パラメータのバリデーションエラー",
		request: map[string]interface{}{
			"title":                "test",
			"discount_amount":      "not_numeric",
			"discount_rate":        "not_numeric",
			"max_discount_amount":  "not_numeric",
			"use_start_at_from":    "not_datetime",
			"use_start_at_to":      "not_datetime",
			"use_end_at_from":      "not_datetime",
			"use_end_at_to":        "not_datetime",
			"public_start_at_from": "not_datetime",
			"public_start_at_to":   "not_datetime",
			"public_end_at_from":   "not_datetime",
			"public_end_at_to":     "not_datetime",
			"is_public":            "not_boolean",
			"is_premium":           "not_boolean",
			"offset":               "not_numeric",
			"limit":                "not_numeric",
		},
		wantStatus: http.StatusBadRequest,
		wantBody: `{
			"code": 400,
			"message": "パラメータが不正です",
			"details": [
				{
					"attribute": "DiscountAmount",
					"message": "DiscountAmountは不正です"
				},
				{
					"attribute": "DiscountRate",
					"message": "DiscountRateは不正です"
				},
				{
					"attribute": "MaxDiscountAmount",
					"message": "MaxDiscountAmountは不正です"
				},
				{
					"attribute": "UseStartAtFrom",
					"message": "UseStartAtFromは不正です"
				},
				{
					"attribute": "UseStartAtTo",
					"message": "UseStartAtToは不正です"
				},
				{
					"attribute": "UseEndAtFrom",
					"message": "UseEndAtFromは不正です"
				},
				{
					"attribute": "UseEndAtTo",
					"message": "UseEndAtToは不正です"
				},
				{
					"attribute": "PublicStartAtFrom",
					"message": "PublicStartAtFromは不正です"
				},
				{
					"attribute": "PublicStartAtTo",
					"message": "PublicStartAtToは不正です"
				},
				{
					"attribute": "PublicEndAtFrom",
					"message": "PublicEndAtFromは不正です"
				},
				{
					"attribute": "PublicEndAtTo",
					"message": "PublicEndAtToは不正です"
				},
				{
					"attribute": "IsPublic",
					"message": "IsPublicは不正です"
				},
				{
					"attribute": "IsPremium",
					"message": "IsPremiumは不正です"
				},
				{
					"attribute": "Offset",
					"message": "Offsetは不正です"
				},
				{
					"attribute": "Limit",
					"message": "Limitは不正です"
				}
			]
		}`,
	},
	{
		tid:  4,
		name: "title で検索",
		request: map[string]interface{}{
			"title": "coupon",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"coupons": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"title": "coupon",
					"remarks": "coupon",
					"discount_amount": 1000,
					"discount_rate": 0,
					"max_discount_amount": 0,
					"use_start_at": "2022-06-01T00:00:00+09:00",
					"use_end_at": "2030-07-02T10:00:00+09:00",
					"public_start_at": "2022-06-01T00:00:00+09:00",
					"public_end_at": "2030-07-02T10:00:00+09:00",
					"is_public": false,
					"is_premium": false,
					"created_at": "2022-06-14T08:19:41+09:00",
					"updated_at": "2022-06-15T10:31:50+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  5,
		name: "discount_amount で範囲検索",
		request: map[string]interface{}{
			"discount_amount": 1000,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"coupons": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"title": "coupon",
					"remarks": "coupon",
					"discount_amount": 1000,
					"discount_rate": 0,
					"max_discount_amount": 0,
					"use_start_at": "2022-06-01T00:00:00+09:00",
					"use_end_at": "2030-07-02T10:00:00+09:00",
					"public_start_at": "2022-06-01T00:00:00+09:00",
					"public_end_at": "2030-07-02T10:00:00+09:00",
					"is_public": false,
					"is_premium": false,
					"created_at": "2022-06-14T08:19:41+09:00",
					"updated_at": "2022-06-15T10:31:50+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  6,
		name: "discount_rate で検索",
		request: map[string]interface{}{
			"discount_rate": 20,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"coupons": [
				{
					"id": "5c3325c1-d539-42d6-b405-2af2f6b99ed9",
					"title": "test",
					"remarks": "test",
					"discount_amount": 0,
					"discount_rate": 20,
					"max_discount_amount": 2000,
					"use_start_at": "2022-07-01T00:00:00+09:00",
					"use_end_at": "2032-07-02T10:00:00+09:00",
					"public_start_at": "2022-07-01T00:00:00+09:00",
					"public_end_at": "2032-07-02T10:00:00+09:00",
					"is_public": true,
					"is_premium": true,
					"created_at": "2022-06-14T08:19:12+09:00",
					"updated_at": "2022-06-14T08:19:12+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  7,
		name: "max_discount_amount で検索",
		request: map[string]interface{}{
			"max_discount_amount": 2000,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"coupons": [
				{
					"id": "5c3325c1-d539-42d6-b405-2af2f6b99ed9",
					"title": "test",
					"remarks": "test",
					"discount_amount": 0,
					"discount_rate": 20,
					"max_discount_amount": 2000,
					"use_start_at": "2022-07-01T00:00:00+09:00",
					"use_end_at": "2032-07-02T10:00:00+09:00",
					"public_start_at": "2022-07-01T00:00:00+09:00",
					"public_end_at": "2032-07-02T10:00:00+09:00",
					"is_public": true,
					"is_premium": true,
					"created_at": "2022-06-14T08:19:12+09:00",
					"updated_at": "2022-06-14T08:19:12+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  8,
		name: "use_start_at_from で範囲検索",
		request: map[string]interface{}{
			"use_start_at_from": "2022-06-25T15:04:05-07:00",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"coupons": [
				{
					"id": "5c3325c1-d539-42d6-b405-2af2f6b99ed9",
					"title": "test",
					"remarks": "test",
					"discount_amount": 0,
					"discount_rate": 20,
					"max_discount_amount": 2000,
					"use_start_at": "2022-07-01T00:00:00+09:00",
					"use_end_at": "2032-07-02T10:00:00+09:00",
					"public_start_at": "2022-07-01T00:00:00+09:00",
					"public_end_at": "2032-07-02T10:00:00+09:00",
					"is_public": true,
					"is_premium": true,
					"created_at": "2022-06-14T08:19:12+09:00",
					"updated_at": "2022-06-14T08:19:12+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  9,
		name: "use_start_at_to で範囲検索",
		request: map[string]interface{}{
			"use_start_at_to": "2022-06-02T15:04:05-07:00",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"coupons": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"title": "coupon",
					"remarks": "coupon",
					"discount_amount": 1000,
					"discount_rate": 0,
					"max_discount_amount": 0,
					"use_start_at": "2022-06-01T00:00:00+09:00",
					"use_end_at": "2030-07-02T10:00:00+09:00",
					"public_start_at": "2022-06-01T00:00:00+09:00",
					"public_end_at": "2030-07-02T10:00:00+09:00",
					"is_public": false,
					"is_premium": false,
					"created_at": "2022-06-14T08:19:41+09:00",
					"updated_at": "2022-06-15T10:31:50+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  10,
		name: "use_end_at_from で範囲検索",
		request: map[string]interface{}{
			"use_end_at_from": "2031-06-02T15:04:05-07:00",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"coupons": [
				{
					"id": "5c3325c1-d539-42d6-b405-2af2f6b99ed9",
					"title": "test",
					"remarks": "test",
					"discount_amount": 0,
					"discount_rate": 20,
					"max_discount_amount": 2000,
					"use_start_at": "2022-07-01T00:00:00+09:00",
					"use_end_at": "2032-07-02T10:00:00+09:00",
					"public_start_at": "2022-07-01T00:00:00+09:00",
					"public_end_at": "2032-07-02T10:00:00+09:00",
					"is_public": true,
					"is_premium": true,
					"created_at": "2022-06-14T08:19:12+09:00",
					"updated_at": "2022-06-14T08:19:12+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  11,
		name: "use_end_at_to で範囲検索",
		request: map[string]interface{}{
			"use_end_at_to": "2031-06-25T15:04:05-07:00",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"coupons": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"title": "coupon",
					"remarks": "coupon",
					"discount_amount": 1000,
					"discount_rate": 0,
					"max_discount_amount": 0,
					"use_start_at": "2022-06-01T00:00:00+09:00",
					"use_end_at": "2030-07-02T10:00:00+09:00",
					"public_start_at": "2022-06-01T00:00:00+09:00",
					"public_end_at": "2030-07-02T10:00:00+09:00",
					"is_public": false,
					"is_premium": false,
					"created_at": "2022-06-14T08:19:41+09:00",
					"updated_at": "2022-06-15T10:31:50+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  12,
		name: "public_start_at_from で範囲検索",
		request: map[string]interface{}{
			"public_start_at_from": "2022-06-25T15:04:05-07:00",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"coupons": [
				{
					"id": "5c3325c1-d539-42d6-b405-2af2f6b99ed9",
					"title": "test",
					"remarks": "test",
					"discount_amount": 0,
					"discount_rate": 20,
					"max_discount_amount": 2000,
					"use_start_at": "2022-07-01T00:00:00+09:00",
					"use_end_at": "2032-07-02T10:00:00+09:00",
					"public_start_at": "2022-07-01T00:00:00+09:00",
					"public_end_at": "2032-07-02T10:00:00+09:00",
					"is_public": true,
					"is_premium": true,
					"created_at": "2022-06-14T08:19:12+09:00",
					"updated_at": "2022-06-14T08:19:12+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  13,
		name: "public_start_at_to で範囲検索",
		request: map[string]interface{}{
			"public_start_at_to": "2022-06-25T15:04:05-07:00",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"coupons": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"title": "coupon",
					"remarks": "coupon",
					"discount_amount": 1000,
					"discount_rate": 0,
					"max_discount_amount": 0,
					"use_start_at": "2022-06-01T00:00:00+09:00",
					"use_end_at": "2030-07-02T10:00:00+09:00",
					"public_start_at": "2022-06-01T00:00:00+09:00",
					"public_end_at": "2030-07-02T10:00:00+09:00",
					"is_public": false,
					"is_premium": false,
					"created_at": "2022-06-14T08:19:41+09:00",
					"updated_at": "2022-06-15T10:31:50+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  14,
		name: "public_end_at_from で範囲検索",
		request: map[string]interface{}{
			"public_end_at_from": "2031-06-02T15:04:05-07:00",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"coupons": [
				{
					"id": "5c3325c1-d539-42d6-b405-2af2f6b99ed9",
					"title": "test",
					"remarks": "test",
					"discount_amount": 0,
					"discount_rate": 20,
					"max_discount_amount": 2000,
					"use_start_at": "2022-07-01T00:00:00+09:00",
					"use_end_at": "2032-07-02T10:00:00+09:00",
					"public_start_at": "2022-07-01T00:00:00+09:00",
					"public_end_at": "2032-07-02T10:00:00+09:00",
					"is_public": true,
					"is_premium": true,
					"created_at": "2022-06-14T08:19:12+09:00",
					"updated_at": "2022-06-14T08:19:12+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  15,
		name: "public_end_at_to で範囲検索",
		request: map[string]interface{}{
			"public_end_at_to": "2031-06-02T15:04:05-07:00",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"coupons": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"title": "coupon",
					"remarks": "coupon",
					"discount_amount": 1000,
					"discount_rate": 0,
					"max_discount_amount": 0,
					"use_start_at": "2022-06-01T00:00:00+09:00",
					"use_end_at": "2030-07-02T10:00:00+09:00",
					"public_start_at": "2022-06-01T00:00:00+09:00",
					"public_end_at": "2030-07-02T10:00:00+09:00",
					"is_public": false,
					"is_premium": false,
					"created_at": "2022-06-14T08:19:41+09:00",
					"updated_at": "2022-06-15T10:31:50+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  16,
		name: "is_public で検索",
		request: map[string]interface{}{
			"is_public": true,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"coupons": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"title": "coupon",
					"remarks": "coupon",
					"discount_amount": 1000,
					"discount_rate": 0,
					"max_discount_amount": 0,
					"use_start_at": "2022-06-01T00:00:00+09:00",
					"use_end_at": "2030-07-02T10:00:00+09:00",
					"public_start_at": "2022-06-01T00:00:00+09:00",
					"public_end_at": "2030-07-02T10:00:00+09:00",
					"is_public": false,
					"is_premium": false,
					"created_at": "2022-06-14T08:19:41+09:00",
					"updated_at": "2022-06-15T10:31:50+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  17,
		name: "is_premium で検索",
		request: map[string]interface{}{
			"is_premium": true,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"coupons": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"title": "coupon",
					"remarks": "coupon",
					"discount_amount": 1000,
					"discount_rate": 0,
					"max_discount_amount": 0,
					"use_start_at": "2022-06-01T00:00:00+09:00",
					"use_end_at": "2030-07-02T10:00:00+09:00",
					"public_start_at": "2022-06-01T00:00:00+09:00",
					"public_end_at": "2030-07-02T10:00:00+09:00",
					"is_public": false,
					"is_premium": false,
					"created_at": "2022-06-14T08:19:41+09:00",
					"updated_at": "2022-06-15T10:31:50+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  18,
		name: "offset 指定で検索",
		request: map[string]interface{}{
			"offset": 1,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"coupons": [
				{
					"id": "5c3325c1-d539-42d6-b405-2af2f6b99ed9",
					"title": "test",
					"remarks": "test",
					"discount_amount": 0,
					"discount_rate": 20,
					"max_discount_amount": 2000,
					"use_start_at": "2022-07-01T00:00:00+09:00",
					"use_end_at": "2032-07-02T10:00:00+09:00",
					"public_start_at": "2022-07-01T00:00:00+09:00",
					"public_end_at": "2032-07-02T10:00:00+09:00",
					"is_public": true,
					"is_premium": true,
					"created_at": "2022-06-14T08:19:12+09:00",
					"updated_at": "2022-06-14T08:19:12+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  19,
		name: "limit で範囲検索",
		request: map[string]interface{}{
			"limit": 1,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"coupons": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"title": "coupon",
					"remarks": "coupon",
					"discount_amount": 1000,
					"discount_rate": 0,
					"max_discount_amount": 0,
					"use_start_at": "2022-06-01T00:00:00+09:00",
					"use_end_at": "2030-07-02T10:00:00+09:00",
					"public_start_at": "2022-06-01T00:00:00+09:00",
					"public_end_at": "2030-07-02T10:00:00+09:00",
					"is_public": false,
					"is_premium": false,
					"created_at": "2022-06-14T08:19:41+09:00",
					"updated_at": "2022-06-15T10:31:50+09:00"
				}
			],
			"total": 1
		}`,
	},
}

func TestGetCoupons(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE coupons")
	require.NoError(t, dbConn.Exec("INSERT INTO coupons (id, title, remarks, discount_amount, discount_rate, max_discount_amount, use_start_at,use_end_at, public_start_at, public_end_at, is_public, is_premium, created_at, updated_at) VALUES ('090e142d-baa3-4039-9d21-cf5a1af39094', 'coupon', 'coupon', 1000, null, null, '2022-06-01 00:00:00','2030-07-02 10:00:00', '2022-06-01 00:00:00', '2030-07-02 10:00:00', 0, 0, '2022-06-14 08:19:41','2022-06-15 10:31:50'),('5c3325c1-d539-42d6-b405-2af2f6b99ed9', 'test', 'test', null, 20, 2000, '2022-07-01 00:00:00','2032-07-02 10:00:00', '2022-07-01 00:00:00', '2032-07-02 10:00:00', 1, 1, '2022-06-14 08:19:12','2022-06-14 08:19:12');").Error)
	r := gin.New()
	zapLogger, err := logger.NewLogger(true)
	require.NoError(t, err)
	r.Use(func(_ *gin.Context) { binding.EnableDecoderUseNumber = true })
	r.Use(middleware.NewTracing())
	r.Use(middleware.NewLogging(zapLogger))
	couponHandler := NewCouponHandler(dbConn, zapLogger)
	r.GET("/coupons", couponHandler.GetAllCoupon)

	for _, tt := range getCouponsTestCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var req *http.Request
			var query string
			for k, v := range tt.request {
				query = query + fmt.Sprintf("&%s=%s", k, url.QueryEscape(fmt.Sprint(v)))
			}
			if query != "" {
				query = "?" + query
			}
			rec := httptest.NewRecorder()
			req = httptest.NewRequest(http.MethodGet, "/coupons"+query, nil)
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}

var getCouponDetailTestCases = []struct {
	tid        int
	name       string
	couponID   string
	wantStatus int
	wantBody   string
}{
	{
		tid:        1,
		name:       "クーポン詳細が正常に取得できること",
		couponID:   couponIDForTest,
		wantStatus: http.StatusOK,
		wantBody: `{
			"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
			"title": "coupon",
			"remarks": "coupon",
			"discount_amount": 1000,
			"discount_rate": 0,
			"max_discount_amount": 0,
			"use_start_at": "2022-06-01T00:00:00+09:00",
			"use_end_at": "2030-07-02T10:00:00+09:00",
			"public_start_at": "2022-06-01T00:00:00+09:00",
			"public_end_at": "2030-07-02T10:00:00+09:00",
			"is_public": false,
			"is_premium": false,
			"created_at": "2022-06-14T08:19:41+09:00",
			"updated_at": "2022-06-15T10:31:50+09:00"
		}`,
	},
	{
		tid:        2,
		name:       "存在しないIDを指定した場合404エラーになること",
		couponID:   "invalid_coupon",
		wantStatus: http.StatusNotFound,
		wantBody:   `{"message": "coupon not found","status": 404,"error": "not_found","causes": null}`,
	},
}

func TestCouponDetail(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE coupons")
	require.NoError(t, dbConn.Exec("INSERT INTO coupons (id, title, remarks, discount_amount, discount_rate, max_discount_amount, use_start_at,use_end_at, public_start_at, public_end_at, is_public, is_premium, created_at, updated_at) VALUES ('090e142d-baa3-4039-9d21-cf5a1af39094', 'coupon', 'coupon', 1000, null, null, '2022-06-01 00:00:00','2030-07-02 10:00:00', '2022-06-01 00:00:00', '2030-07-02 10:00:00', 0, 0, '2022-06-14 08:19:41','2022-06-15 10:31:50'),('5c3325c1-d539-42d6-b405-2af2f6b99ed9', 'test', 'test', null, 20, 2000, '2022-07-01 00:00:00','2032-07-02 10:00:00', '2022-07-01 00:00:00', '2032-07-02 10:00:00', 1, 1, '2022-06-14 08:19:12','2022-06-14 08:19:12');").Error)
	r := gin.New()
	zapLogger, err := logger.NewLogger(true)
	require.NoError(t, err)
	r.Use(func(_ *gin.Context) { binding.EnableDecoderUseNumber = true })
	r.Use(middleware.NewTracing())
	r.Use(middleware.NewLogging(zapLogger))
	couponHandler := NewCouponHandler(dbConn, zapLogger)
	r.GET("/coupons/:id", couponHandler.GetCouponDetail)

	for _, tt := range getCouponDetailTestCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var req *http.Request
			rec := httptest.NewRecorder()
			req = httptest.NewRequest(http.MethodGet, "/coupons/"+tt.couponID, nil)
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}
