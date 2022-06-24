package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/AI1411/golang-admin-api/middleware"
	logger "github.com/AI1411/golang-admin-api/server"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/AI1411/golang-admin-api/db"
)

const (
	orderIDForTest     = "090e142d-baa3-4039-9d21-cf5a1af39094"
	userIDForOrderTest = "7dc41179-824e-4b8a-b894-2082ca5eac5b"
)

var getOrdersTestCases = []struct {
	tid        int
	name       string
	request    map[string]interface{}
	wantStatus int
	wantBody   string
}{
	{
		tid:        1,
		name:       "注文一覧が正常に取得できること",
		request:    map[string]interface{}{},
		wantStatus: http.StatusOK,
		wantBody: `{
			"orders": [
				{
					"id": "5c3325c1-d539-42d6-b405-2af2f6b99ed9",
					"user_id": "db64d2b0-76b0-41f5-8519-89d03a26dde3",
					"quantity": 2,
					"total_price": 200,
					"order_status": "waiting",
					"remarks": "remarks",
					"order_details": [
						{
							"id": "27f1c6ce-7588-4300-9014-e6649af06319",
							"order_id": "5c3325c1-d539-42d6-b405-2af2f6b99ed9",
							"product_id": "66925ce2-47ee-4dfb-b974-f0ba3cd5c177",
							"quantity": 1,
							"order_detail_status": "new",
							"price": 1000
						},
						{
							"id": "2a839d06-c61b-49f8-bec1-a9a604d11db0",
							"order_id": "5c3325c1-d539-42d6-b405-2af2f6b99ed9",
							"product_id": "66925ce2-47ee-4dfb-b974-f0ba3cd5c177",
							"quantity": 1,
							"order_detail_status": "new",
							"price": 2000
						}
					],
					"created_at": "2022-06-11T19:59:01+09:00",
					"updated_at": "2022-06-11T19:59:01+09:00"
				},
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"user_id": "7dc41179-824e-4b8a-b894-2082ca5eac5b",
					"quantity": 1,
					"total_price": 100,
					"order_status": "new",
					"remarks": "test",
					"order_details": [
						{
							"id": "218c51c0-904e-4743-a2ae-94f0e34a0d6f",
							"order_id": "090e142d-baa3-4039-9d21-cf5a1af39094",
							"product_id": "66925ce2-47ee-4dfb-b974-f0ba3cd5c178",
							"quantity": 1,
							"order_detail_status": "new",
							"price": 100
						},
						{
							"id": "23c66d26-4432-4f7f-9a0d-2642731a28cc",
							"order_id": "090e142d-baa3-4039-9d21-cf5a1af39094",
							"product_id": "1583cc19-bbfa-405a-affb-9f01953f5b6d",
							"quantity": 1,
							"order_detail_status": "new",
							"price": 200
						}
					],
					"created_at": "2022-06-11T10:36:43+09:00",
					"updated_at": "2022-06-11T10:36:43+09:00"
				}
			],
			"total": 2
		}`,
	},
	{
		tid:  2,
		name: "検索結果0件",
		request: map[string]interface{}{
			"user_id": userIDForTest,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"orders": [],
    		"total": 0
		}`,
	},
	{
		tid:  3,
		name: "パラメータのバリデーションエラー",
		request: map[string]interface{}{
			"user_id":      "test",
			"quantity":     "not_numeric",
			"total_price":  "not_numeric",
			"order_status": "test",
			"offset":       "not_numeric",
			"limit":        "not_numeric",
		},
		wantStatus: http.StatusBadRequest,
		wantBody: `{
			"code": 400,
			"message": "パラメータが不正です",
			"details": [
				{
					"attribute": "UserID",
					"message": "ユーザーIDは不正です"
				},
				{
					"attribute": "Quantity",
					"message": "数量は不正です"
				},
				{
					"attribute": "TotalPrice",
					"message": "合計金額は不正です"
				},
				{
					"attribute": "OrderStatus",
					"message": "注文ステータスは不正です"
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
		name: "user_id で検索",
		request: map[string]interface{}{
			"user_id": userIDForOrderTest,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"orders": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"user_id": "7dc41179-824e-4b8a-b894-2082ca5eac5b",
					"quantity": 1,
					"total_price": 100,
					"order_status": "new",
					"remarks": "test",
					"order_details": [
						{
							"id": "218c51c0-904e-4743-a2ae-94f0e34a0d6f",
							"order_id": "090e142d-baa3-4039-9d21-cf5a1af39094",
							"product_id": "66925ce2-47ee-4dfb-b974-f0ba3cd5c178",
							"quantity": 1,
							"order_detail_status": "new",
							"price": 100
						},
						{
							"id": "23c66d26-4432-4f7f-9a0d-2642731a28cc",
							"order_id": "090e142d-baa3-4039-9d21-cf5a1af39094",
							"product_id": "1583cc19-bbfa-405a-affb-9f01953f5b6d",
							"quantity": 1,
							"order_detail_status": "new",
							"price": 200
						}
					],
					"created_at": "2022-06-11T10:36:43+09:00",
					"updated_at": "2022-06-11T10:36:43+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  5,
		name: "quantity で範囲検索",
		request: map[string]interface{}{
			"quantity": 1,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"orders": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"user_id": "7dc41179-824e-4b8a-b894-2082ca5eac5b",
					"quantity": 1,
					"total_price": 100,
					"order_status": "new",
					"remarks": "test",
					"order_details": [
						{
							"id": "218c51c0-904e-4743-a2ae-94f0e34a0d6f",
							"order_id": "090e142d-baa3-4039-9d21-cf5a1af39094",
							"product_id": "66925ce2-47ee-4dfb-b974-f0ba3cd5c178",
							"quantity": 1,
							"order_detail_status": "new",
							"price": 100
						},
						{
							"id": "23c66d26-4432-4f7f-9a0d-2642731a28cc",
							"order_id": "090e142d-baa3-4039-9d21-cf5a1af39094",
							"product_id": "1583cc19-bbfa-405a-affb-9f01953f5b6d",
							"quantity": 1,
							"order_detail_status": "new",
							"price": 200
						}
					],
					"created_at": "2022-06-11T10:36:43+09:00",
					"updated_at": "2022-06-11T10:36:43+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  6,
		name: "total_price で検索",
		request: map[string]interface{}{
			"total_price": 100,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"orders": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"user_id": "7dc41179-824e-4b8a-b894-2082ca5eac5b",
					"quantity": 1,
					"total_price": 100,
					"order_status": "new",
					"remarks": "test",
					"order_details": [
						{
							"id": "218c51c0-904e-4743-a2ae-94f0e34a0d6f",
							"order_id": "090e142d-baa3-4039-9d21-cf5a1af39094",
							"product_id": "66925ce2-47ee-4dfb-b974-f0ba3cd5c178",
							"quantity": 1,
							"order_detail_status": "new",
							"price": 100
						},
						{
							"id": "23c66d26-4432-4f7f-9a0d-2642731a28cc",
							"order_id": "090e142d-baa3-4039-9d21-cf5a1af39094",
							"product_id": "1583cc19-bbfa-405a-affb-9f01953f5b6d",
							"quantity": 1,
							"order_detail_status": "new",
							"price": 200
						}
					],
					"created_at": "2022-06-11T10:36:43+09:00",
					"updated_at": "2022-06-11T10:36:43+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  7,
		name: "order_status で検索",
		request: map[string]interface{}{
			"order_status": "new",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"orders": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"user_id": "7dc41179-824e-4b8a-b894-2082ca5eac5b",
					"quantity": 1,
					"total_price": 100,
					"order_status": "new",
					"remarks": "test",
					"order_details": [
						{
							"id": "218c51c0-904e-4743-a2ae-94f0e34a0d6f",
							"order_id": "090e142d-baa3-4039-9d21-cf5a1af39094",
							"product_id": "66925ce2-47ee-4dfb-b974-f0ba3cd5c178",
							"quantity": 1,
							"order_detail_status": "new",
							"price": 100
						},
						{
							"id": "23c66d26-4432-4f7f-9a0d-2642731a28cc",
							"order_id": "090e142d-baa3-4039-9d21-cf5a1af39094",
							"product_id": "1583cc19-bbfa-405a-affb-9f01953f5b6d",
							"quantity": 1,
							"order_detail_status": "new",
							"price": 200
						}
					],
					"created_at": "2022-06-11T10:36:43+09:00",
					"updated_at": "2022-06-11T10:36:43+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  8,
		name: "offset 指定で検索",
		request: map[string]interface{}{
			"offset": 1,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"orders": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"user_id": "7dc41179-824e-4b8a-b894-2082ca5eac5b",
					"quantity": 1,
					"total_price": 100,
					"order_status": "new",
					"remarks": "test",
					"order_details": [
						{
							"id": "218c51c0-904e-4743-a2ae-94f0e34a0d6f",
							"order_id": "090e142d-baa3-4039-9d21-cf5a1af39094",
							"product_id": "66925ce2-47ee-4dfb-b974-f0ba3cd5c178",
							"quantity": 1,
							"order_detail_status": "new",
							"price": 100
						},
						{
							"id": "23c66d26-4432-4f7f-9a0d-2642731a28cc",
							"order_id": "090e142d-baa3-4039-9d21-cf5a1af39094",
							"product_id": "1583cc19-bbfa-405a-affb-9f01953f5b6d",
							"quantity": 1,
							"order_detail_status": "new",
							"price": 200
						}
					],
					"created_at": "2022-06-11T10:36:43+09:00",
					"updated_at": "2022-06-11T10:36:43+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  9,
		name: "limit 指定で検索",
		request: map[string]interface{}{
			"limit": 1,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"orders": [
				{
					"id": "5c3325c1-d539-42d6-b405-2af2f6b99ed9",
					"user_id": "db64d2b0-76b0-41f5-8519-89d03a26dde3",
					"quantity": 2,
					"total_price": 200,
					"order_status": "waiting",
					"remarks": "remarks",
					"order_details": [
						{
							"id": "27f1c6ce-7588-4300-9014-e6649af06319",
							"order_id": "5c3325c1-d539-42d6-b405-2af2f6b99ed9",
							"product_id": "66925ce2-47ee-4dfb-b974-f0ba3cd5c177",
							"quantity": 1,
							"order_detail_status": "new",
							"price": 1000
						},
						{
							"id": "2a839d06-c61b-49f8-bec1-a9a604d11db0",
							"order_id": "5c3325c1-d539-42d6-b405-2af2f6b99ed9",
							"product_id": "66925ce2-47ee-4dfb-b974-f0ba3cd5c177",
							"quantity": 1,
							"order_detail_status": "new",
							"price": 2000
						}
					],
					"created_at": "2022-06-11T19:59:01+09:00",
					"updated_at": "2022-06-11T19:59:01+09:00"
				}
			],
			"total": 1
		}`,
	},
}

func TestGetOrders(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE orders")
	dbConn.Exec("TRUNCATE TABLE order_details")
	require.NoError(t, dbConn.Exec("INSERT INTO orders (id, user_id, quantity, total_price, order_status, remarks, created_at, updated_at)VALUES ('090e142d-baa3-4039-9d21-cf5a1af39094', '7dc41179-824e-4b8a-b894-2082ca5eac5b', '1', 100, 'new', 'test','2022-06-11 10:36:43', '2022-06-11 10:36:43'), ('5c3325c1-d539-42d6-b405-2af2f6b99ed9', 'db64d2b0-76b0-41f5-8519-89d03a26dde3', '2', 200, 'waiting', 'remarks', '2022-06-11 19:59:01', '2022-06-11 19:59:01');").Error)
	require.NoError(t, dbConn.Exec("INSERT INTO order_details(id,order_id,product_id,quantity,price,order_detail_status,created_at,updated_at)VALUES('218c51c0-904e-4743-a2ae-94f0e34a0d6f','090e142d-baa3-4039-9d21-cf5a1af39094','66925ce2-47ee-4dfb-b974-f0ba3cd5c178','1',100,'new','2022-06-01 16:10:15','2022-06-01 16:10:15'),('23c66d26-4432-4f7f-9a0d-2642731a28cc','090e142d-baa3-4039-9d21-cf5a1af39094','1583cc19-bbfa-405a-affb-9f01953f5b6d','1',200,'new','2022-06-11 14:50:51','2022-06-11 14:50:51'),('27f1c6ce-7588-4300-9014-e6649af06319','5c3325c1-d539-42d6-b405-2af2f6b99ed9','66925ce2-47ee-4dfb-b974-f0ba3cd5c177','1',1000,'new','2022-06-05 15:33:46','2022-06-05 15:33:46'),('2a839d06-c61b-49f8-bec1-a9a604d11db0','5c3325c1-d539-42d6-b405-2af2f6b99ed9','66925ce2-47ee-4dfb-b974-f0ba3cd5c177','1',2000,'new','2022-06-01 16:10:15','2022-06-01 16:10:15');").Error)
	r := gin.New()
	zapLogger, err := logger.NewLogger(true)
	require.NoError(t, err)
	r.Use(func(_ *gin.Context) { binding.EnableDecoderUseNumber = true })
	r.Use(middleware.NewTracing())
	r.Use(middleware.NewLogging(zapLogger))
	orderHandler := NewOrderHandler(dbConn, zapLogger)
	r.GET("/orders", orderHandler.GetOrders)

	for _, tt := range getOrdersTestCases {
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
			req = httptest.NewRequest(http.MethodGet, "/orders"+query, nil)
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}

var getOrderDetailTestCases = []struct {
	tid        int
	name       string
	orderID    string
	wantStatus int
	wantBody   string
}{
	{
		tid:        1,
		name:       "注文詳細が正常に取得できること",
		orderID:    orderIDForTest,
		wantStatus: http.StatusOK,
		wantBody: `{
			"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
			"user_id": "7dc41179-824e-4b8a-b894-2082ca5eac5b",
			"quantity": 1,
			"total_price": 100,
			"order_status": "new",
			"remarks": "test",
			"order_details": [
				{
					"id": "218c51c0-904e-4743-a2ae-94f0e34a0d6f",
					"order_id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"product_id": "66925ce2-47ee-4dfb-b974-f0ba3cd5c178",
					"quantity": 1,
					"order_detail_status": "new",
					"price": 100
				},
				{
					"id": "23c66d26-4432-4f7f-9a0d-2642731a28cc",
					"order_id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"product_id": "1583cc19-bbfa-405a-affb-9f01953f5b6d",
					"quantity": 1,
					"order_detail_status": "new",
					"price": 200
				}
			],
			"created_at": "2022-06-11T10:36:43+09:00",
			"updated_at": "2022-06-11T10:36:43+09:00"
		}`,
	},
	{
		tid:        2,
		name:       "存在しないIDを指定した場合404エラーになること",
		orderID:    "invalid_order",
		wantStatus: http.StatusNotFound,
		wantBody:   `{"message": "order not found","status": 404,"error": "not_found","causes": null}`,
	},
}

func TestOrderDetail(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE orders")
	dbConn.Exec("TRUNCATE TABLE order_details")
	require.NoError(t, dbConn.Exec("INSERT INTO orders (id, user_id, quantity, total_price, order_status, remarks, created_at, updated_at)VALUES ('090e142d-baa3-4039-9d21-cf5a1af39094', '7dc41179-824e-4b8a-b894-2082ca5eac5b', '1', 100, 'new', 'test','2022-06-11 10:36:43', '2022-06-11 10:36:43'), ('5c3325c1-d539-42d6-b405-2af2f6b99ed9', 'db64d2b0-76b0-41f5-8519-89d03a26dde3', '2', 200, 'waiting', 'remarks', '2022-06-11 19:59:01', '2022-06-11 19:59:01');").Error)
	require.NoError(t, dbConn.Exec("INSERT INTO order_details(id,order_id,product_id,quantity,price,order_detail_status,created_at,updated_at)VALUES('218c51c0-904e-4743-a2ae-94f0e34a0d6f','090e142d-baa3-4039-9d21-cf5a1af39094','66925ce2-47ee-4dfb-b974-f0ba3cd5c178','1',100,'new','2022-06-01 16:10:15','2022-06-01 16:10:15'),('23c66d26-4432-4f7f-9a0d-2642731a28cc','090e142d-baa3-4039-9d21-cf5a1af39094','1583cc19-bbfa-405a-affb-9f01953f5b6d','1',200,'new','2022-06-11 14:50:51','2022-06-11 14:50:51'),('27f1c6ce-7588-4300-9014-e6649af06319','5c3325c1-d539-42d6-b405-2af2f6b99ed9','66925ce2-47ee-4dfb-b974-f0ba3cd5c177','1',1000,'new','2022-06-05 15:33:46','2022-06-05 15:33:46'),('2a839d06-c61b-49f8-bec1-a9a604d11db0','5c3325c1-d539-42d6-b405-2af2f6b99ed9','66925ce2-47ee-4dfb-b974-f0ba3cd5c177','1',2000,'new','2022-06-01 16:10:15','2022-06-01 16:10:15');").Error)
	r := gin.New()
	zapLogger, err := logger.NewLogger(true)
	require.NoError(t, err)
	r.Use(func(_ *gin.Context) { binding.EnableDecoderUseNumber = true })
	r.Use(middleware.NewTracing())
	r.Use(middleware.NewLogging(zapLogger))
	orderHandler := NewOrderHandler(dbConn, zapLogger)
	r.GET("/orders/:id", orderHandler.GetOrder)

	for _, tt := range getOrderDetailTestCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var req *http.Request
			rec := httptest.NewRecorder()
			req = httptest.NewRequest(http.MethodGet, "/orders/"+tt.orderID, nil)
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}

var deleteOrderTestCases = []struct {
	tid        int
	name       string
	orderID    string
	request    map[string]interface{}
	wantStatus int
	wantBody   string
}{
	{
		tid:        1,
		name:       "注文が正常に削除できること",
		orderID:    orderIDForTest,
		wantStatus: http.StatusNoContent,
	},
	{
		tid:        2,
		name:       "削除できる注文がない場合は404エラー",
		orderID:    "10",
		wantStatus: http.StatusNotFound,
	},
}

func TestDeleteOrder(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE orders")
	dbConn.Exec("TRUNCATE TABLE order_details")
	require.NoError(t, dbConn.Exec("INSERT INTO orders (id, user_id, quantity, total_price, order_status, remarks, created_at, updated_at)VALUES ('090e142d-baa3-4039-9d21-cf5a1af39094', '7dc41179-824e-4b8a-b894-2082ca5eac5b', '1', 100, 'new', 'test','2022-06-11 10:36:43', '2022-06-11 10:36:43'), ('5c3325c1-d539-42d6-b405-2af2f6b99ed9', 'db64d2b0-76b0-41f5-8519-89d03a26dde3', '2', 200, 'waiting', 'remarks', '2022-06-11 19:59:01', '2022-06-11 19:59:01');").Error)
	require.NoError(t, dbConn.Exec("INSERT INTO order_details(id,order_id,product_id,quantity,price,order_detail_status,created_at,updated_at)VALUES('218c51c0-904e-4743-a2ae-94f0e34a0d6f','090e142d-baa3-4039-9d21-cf5a1af39094','66925ce2-47ee-4dfb-b974-f0ba3cd5c178','1',100,'new','2022-06-01 16:10:15','2022-06-01 16:10:15'),('23c66d26-4432-4f7f-9a0d-2642731a28cc','090e142d-baa3-4039-9d21-cf5a1af39094','1583cc19-bbfa-405a-affb-9f01953f5b6d','1',200,'new','2022-06-11 14:50:51','2022-06-11 14:50:51'),('27f1c6ce-7588-4300-9014-e6649af06319','5c3325c1-d539-42d6-b405-2af2f6b99ed9','66925ce2-47ee-4dfb-b974-f0ba3cd5c177','1',1000,'new','2022-06-05 15:33:46','2022-06-05 15:33:46'),('2a839d06-c61b-49f8-bec1-a9a604d11db0','5c3325c1-d539-42d6-b405-2af2f6b99ed9','66925ce2-47ee-4dfb-b974-f0ba3cd5c177','1',2000,'new','2022-06-01 16:10:15','2022-06-01 16:10:15');").Error)
	r := gin.New()
	zapLogger, err := logger.NewLogger(true)
	require.NoError(t, err)
	r.Use(func(_ *gin.Context) { binding.EnableDecoderUseNumber = true })
	r.Use(middleware.NewTracing())
	r.Use(middleware.NewLogging(zapLogger))
	orderHandler := NewOrderHandler(dbConn, zapLogger)
	r.DELETE("/orders/:id", orderHandler.DeleteOrder)

	for _, tt := range deleteOrderTestCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var req *http.Request
			rec := httptest.NewRecorder()
			jsonStr, _ := json.Marshal(tt.request)
			req = httptest.NewRequest(http.MethodDelete, "/orders/"+tt.orderID, bytes.NewBuffer(jsonStr))
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, rec.Body.String())
			}
		})
	}
}
