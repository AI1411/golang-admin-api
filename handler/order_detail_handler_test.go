package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/AI1411/golang-admin-api/db"
)

const orderDetailIDForTest = "218c51c0-904e-4743-a2ae-94f0e34a0d6f"

var getOrderDetailDetailTestCases = []struct {
	tid           int
	name          string
	orderDetailID string
	wantStatus    int
	wantBody      string
}{
	{
		tid:           1,
		name:          "マイルストーン詳細が正常に取得できること",
		orderDetailID: orderDetailIDForTest,
		wantStatus:    http.StatusOK,
		wantBody: `{
			"id": "218c51c0-904e-4743-a2ae-94f0e34a0d6f",
			"order_id": "090e142d-baa3-4039-9d21-cf5a1af39094",
			"product_id": "66925ce2-47ee-4dfb-b974-f0ba3cd5c178",
			"quantity": 1,
			"order_detail_status": "new",
			"price": 100
		}`,
	},
	{
		tid:           2,
		name:          "存在しないIDを指定した場合404エラーになること",
		orderDetailID: "invalid_orderDetail",
		wantStatus:    http.StatusNotFound,
		wantBody:      `{"message": "order detail not found","status": 404,"error": "not_found","causes": null}`,
	},
}

func TestOrderDetailDetail(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE orders")
	dbConn.Exec("TRUNCATE TABLE order_details")
	require.NoError(t, dbConn.Exec("INSERT INTO orders (id, user_id, quantity, total_price, order_status, remarks, created_at, updated_at)VALUES ('090e142d-baa3-4039-9d21-cf5a1af39094', '7dc41179-824e-4b8a-b894-2082ca5eac5b', '1', 100, 'new', 'test','2022-06-11 10:36:43', '2022-06-11 10:36:43'), ('5c3325c1-d539-42d6-b405-2af2f6b99ed9', 'db64d2b0-76b0-41f5-8519-89d03a26dde3', '2', 200, 'waiting', 'remarks', '2022-06-11 19:59:01', '2022-06-11 19:59:01');").Error)
	require.NoError(t, dbConn.Exec("INSERT INTO order_details(id,order_id,product_id,quantity,price,order_detail_status,created_at,updated_at)VALUES('218c51c0-904e-4743-a2ae-94f0e34a0d6f','090e142d-baa3-4039-9d21-cf5a1af39094','66925ce2-47ee-4dfb-b974-f0ba3cd5c178','1',100,'new','2022-06-01 16:10:15','2022-06-01 16:10:15'),('23c66d26-4432-4f7f-9a0d-2642731a28cc','090e142d-baa3-4039-9d21-cf5a1af39094','1583cc19-bbfa-405a-affb-9f01953f5b6d','1',200,'new','2022-06-11 14:50:51','2022-06-11 14:50:51'),('27f1c6ce-7588-4300-9014-e6649af06319','5c3325c1-d539-42d6-b405-2af2f6b99ed9','66925ce2-47ee-4dfb-b974-f0ba3cd5c177','1',1000,'new','2022-06-05 15:33:46','2022-06-05 15:33:46'),('2a839d06-c61b-49f8-bec1-a9a604d11db0','5c3325c1-d539-42d6-b405-2af2f6b99ed9','66925ce2-47ee-4dfb-b974-f0ba3cd5c177','1',2000,'new','2022-06-01 16:10:15','2022-06-01 16:10:15');").Error)
	r := gin.New()
	orderDetailHandler := NewOrderDetailHandler(dbConn)
	r.GET("/orderDetails/:id", orderDetailHandler.GetOrderDetail)

	for _, tt := range getOrderDetailDetailTestCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var req *http.Request
			rec := httptest.NewRecorder()
			req = httptest.NewRequest(http.MethodGet, "/orderDetails/"+tt.orderDetailID, nil)
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}

var deleteOrderDetailTestCases = []struct {
	tid           int
	name          string
	orderDetailID string
	request       map[string]interface{}
	wantStatus    int
	wantBody      string
}{
	{
		tid:           1,
		name:          "マイルストーンが正常に削除できること",
		orderDetailID: orderDetailIDForTest,
		wantStatus:    http.StatusNoContent,
	},
	{
		tid:           2,
		name:          "削除できるマイルストーンがない場合は404エラー",
		orderDetailID: "10",
		wantStatus:    http.StatusNotFound,
	},
}

func TestDeleteOrderDetail(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE orders")
	dbConn.Exec("TRUNCATE TABLE order_details")
	require.NoError(t, dbConn.Exec("INSERT INTO orders (id, user_id, quantity, total_price, order_status, remarks, created_at, updated_at)VALUES ('090e142d-baa3-4039-9d21-cf5a1af39094', '7dc41179-824e-4b8a-b894-2082ca5eac5b', '1', 100, 'new', 'test','2022-06-11 10:36:43', '2022-06-11 10:36:43'), ('5c3325c1-d539-42d6-b405-2af2f6b99ed9', 'db64d2b0-76b0-41f5-8519-89d03a26dde3', '2', 200, 'waiting', 'remarks', '2022-06-11 19:59:01', '2022-06-11 19:59:01');").Error)
	require.NoError(t, dbConn.Exec("INSERT INTO order_details(id,order_id,product_id,quantity,price,order_detail_status,created_at,updated_at)VALUES('218c51c0-904e-4743-a2ae-94f0e34a0d6f','090e142d-baa3-4039-9d21-cf5a1af39094','66925ce2-47ee-4dfb-b974-f0ba3cd5c178','1',100,'new','2022-06-01 16:10:15','2022-06-01 16:10:15'),('23c66d26-4432-4f7f-9a0d-2642731a28cc','090e142d-baa3-4039-9d21-cf5a1af39094','1583cc19-bbfa-405a-affb-9f01953f5b6d','1',200,'new','2022-06-11 14:50:51','2022-06-11 14:50:51'),('27f1c6ce-7588-4300-9014-e6649af06319','5c3325c1-d539-42d6-b405-2af2f6b99ed9','66925ce2-47ee-4dfb-b974-f0ba3cd5c177','1',1000,'new','2022-06-05 15:33:46','2022-06-05 15:33:46'),('2a839d06-c61b-49f8-bec1-a9a604d11db0','5c3325c1-d539-42d6-b405-2af2f6b99ed9','66925ce2-47ee-4dfb-b974-f0ba3cd5c177','1',2000,'new','2022-06-01 16:10:15','2022-06-01 16:10:15');").Error)
	r := gin.New()
	orderDetailHandler := NewOrderDetailHandler(dbConn)
	r.DELETE("/orderDetails/:id", orderDetailHandler.DeleteOrderDetail)

	for _, tt := range deleteOrderDetailTestCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var req *http.Request
			rec := httptest.NewRecorder()
			jsonStr, _ := json.Marshal(tt.request)
			req = httptest.NewRequest(http.MethodDelete, "/orderDetails/"+tt.orderDetailID, bytes.NewBuffer(jsonStr))
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, rec.Body.String())
			}
		})
	}
}
