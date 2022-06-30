package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/AI1411/golang-admin-api/middleware"
	logger "github.com/AI1411/golang-admin-api/server"
	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/require"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/AI1411/golang-admin-api/db"
	"github.com/AI1411/golang-admin-api/models/mock_model"
)

const productIDForTest = "090e142d-baa3-4039-9d21-cf5a1af39094"

var getProductsTestCases = []struct {
	tid        int
	name       string
	request    map[string]interface{}
	wantStatus int
	wantBody   string
}{
	{
		tid:        1,
		name:       "商品一覧が正常に取得できること",
		request:    map[string]interface{}{},
		wantStatus: http.StatusOK,
		wantBody: `{
			"products": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"product_name": "1",
					"price": 100,
					"remarks": "1",
					"quantity": 1
				},
				{
					"id": "5c3325c1-d539-42d6-b405-2af2f6b99ed9",
					"product_name": "2",
					"price": 1000,
					"remarks": "2",
					"quantity": 10
				}
			],
			"total": 2
		}`,
	},
	{
		tid:  2,
		name: "検索結果0件",
		request: map[string]interface{}{
			"product_name": "failed",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"products": [],
    		"total": 0
		}`,
	},
	{
		tid:  3,
		name: "パラメータのバリデーションエラー",
		request: map[string]interface{}{
			"product_name": strings.Repeat("a", 65),
			"price_from":   "t",
			"price_to":     "t",
			"remarks":      strings.Repeat("a", 256),
			"quantity":     0,
			"offset":       "not_numeric",
			"limit":        "not_numeric",
		},
		wantStatus: http.StatusBadRequest,
		wantBody: `{
			"code": 400,
			"message": "パラメータが不正です",
			"details": [
				{
					"attribute": "ProductName",
					"message": "ProductNameは不正です"
				},
				{
					"attribute": "PriceFrom",
					"message": "PriceFromは不正です"
				},
				{
					"attribute": "PriceTo",
					"message": "PriceToは不正です"
				},
				{
					"attribute": "Remarks",
					"message": "備考は不正です"
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
		name: "product_name で検索",
		request: map[string]interface{}{
			"product_name": "1",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"products": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"product_name": "1",
					"price": 100,
					"remarks": "1",
					"quantity": 1
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  5,
		name: "price_from で範囲検索",
		request: map[string]interface{}{
			"price_from": 200,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"products": [
				{
					"id": "5c3325c1-d539-42d6-b405-2af2f6b99ed9",
					"product_name": "2",
					"price": 1000,
					"remarks": "2",
					"quantity": 10
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  6,
		name: "price_to で範囲検索",
		request: map[string]interface{}{
			"price_to": 101,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"products": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"product_name": "1",
					"price": 100,
					"remarks": "1",
					"quantity": 1
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  7,
		name: "remarks で検索",
		request: map[string]interface{}{
			"remarks": "1",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"products": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"product_name": "1",
					"price": 100,
					"remarks": "1",
					"quantity": 1
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  8,
		name: "quantity で検索",
		request: map[string]interface{}{
			"quantity": 1,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"products": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"product_name": "1",
					"price": 100,
					"remarks": "1",
					"quantity": 1
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
			"products": [
				{
					"id": "5c3325c1-d539-42d6-b405-2af2f6b99ed9",
					"product_name": "2",
					"price": 1000,
					"remarks": "2",
					"quantity": 10
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
			"products": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"product_name": "1",
					"price": 100,
					"remarks": "1",
					"quantity": 1
				}
			],
			"total": 1
		}`,
	},
}

func TestGetProducts(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE products")
	dbConn.Exec("insert into products (id, product_name, price,remarks,quantity, created_at, updated_at)values('090e142d-baa3-4039-9d21-cf5a1af39094','1',100,'1',1,'2022-06-20 22:14:22','2022-06-20 22:14:22'),('5c3325c1-d539-42d6-b405-2af2f6b99ed9','2',1000,'2',10, '2022-06-20 22:14:23', '2022-06-20 22:14:23');")
	r := gin.New()
	zapLogger, err := logger.NewLoggerForTest(true)
	require.NoError(t, err)
	r.Use(func(_ *gin.Context) { binding.EnableDecoderUseNumber = true })
	r.Use(middleware.NewTracing())
	r.Use(middleware.NewLogging(zapLogger))
	productHandler := NewProductHandler(dbConn, nil, zapLogger)
	r.GET("/products", productHandler.GetAllProduct)

	for _, tt := range getProductsTestCases {
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
			req = httptest.NewRequest(http.MethodGet, "/products"+query, nil)
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}

var getProductDetailTestCases = []struct {
	tid        int
	name       string
	productID  string
	wantStatus int
	wantBody   string
}{
	{
		tid:        1,
		name:       "商品詳細が正常に取得できること",
		productID:  productIDForTest,
		wantStatus: http.StatusOK,
		wantBody: `{
			"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
			"product_name": "1",
			"price": 100,
			"remarks": "1",
			"quantity": 1
		}`,
	},
	{
		tid:        2,
		name:       "存在しないIDを指定した場合404エラーになること",
		productID:  "invalid_product",
		wantStatus: http.StatusNotFound,
		wantBody:   `{"message": "product not found","status": 404,"error": "not_found","causes": null}`,
	},
}

func TestProductDetail(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE products")
	dbConn.Exec("insert into products (id, product_name, price,remarks,quantity, created_at, updated_at)values('090e142d-baa3-4039-9d21-cf5a1af39094','1',100,'1',1,'2022-06-20 22:14:22','2022-06-20 22:14:22'),('5c3325c1-d539-42d6-b405-2af2f6b99ed9','2',1000,'2',10, '2022-06-20 22:14:23', '2022-06-20 22:14:23');")
	r := gin.New()
	zapLogger, err := logger.NewLoggerForTest(true)
	require.NoError(t, err)
	r.Use(func(_ *gin.Context) { binding.EnableDecoderUseNumber = true })
	r.Use(middleware.NewTracing())
	r.Use(middleware.NewLogging(zapLogger))
	productHandler := NewProductHandler(dbConn, nil, zapLogger)
	r.GET("/products/:id", productHandler.GetProductDetail)

	for _, tt := range getProductDetailTestCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var req *http.Request
			rec := httptest.NewRecorder()
			req = httptest.NewRequest(http.MethodGet, "/products/"+tt.productID, nil)
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}

var createProductTestCases = []struct {
	tid        int
	name       string
	request    map[string]interface{}
	wantStatus int
	wantBody   string
}{
	{
		tid:  1,
		name: "プロジェクトが正常に作成できること",
		request: map[string]interface{}{
			"product_name": "test",
			"price":        1,
			"remarks":      "update",
			"quantity":     1,
		},
		wantStatus: http.StatusCreated,
		wantBody: `{
			"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
			"product_name": "test",
			"price": 1,
			"remarks": "update",
			"quantity": 1
		}`,
	},
	{
		tid:  2,
		name: "バリデーションエラー",
		request: map[string]interface{}{
			"product_name": strings.Repeat("a", 65),
			"remarks":      strings.Repeat("a", 256),
		},
		wantStatus: http.StatusBadRequest,
		wantBody: `{
			"code": 400,
			"message": "パラメータが不正です",
			"details": [
				{
					"attribute": "ProductName",
					"message": "ProductNameは不正です"
				},
				{
					"attribute": "Price",
					"message": "Priceは必須です"
				},
				{
					"attribute": "Remarks",
					"message": "備考は不正です"
				},
				{
					"attribute": "Quantity",
					"message": "数量は必須です"
				}
			]
		}`,
	},
}

func TestCreateProduct(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE products")
	r := gin.New()
	zapLogger, err := logger.NewLoggerForTest(true)
	require.NoError(t, err)
	r.Use(func(_ *gin.Context) { binding.EnableDecoderUseNumber = true })
	r.Use(middleware.NewTracing())
	r.Use(middleware.NewLogging(zapLogger))
	mockCtrl := gomock.NewController(t)
	uuidGen := mock_models.NewMockUUIDGenerator(mockCtrl)
	uuidGen.EXPECT().GenerateUUID().Return(productIDForTest)

	productHandler := NewProductHandler(dbConn, uuidGen, zapLogger)
	r.POST("/products", productHandler.CreateProduct)

	for _, tt := range createProductTestCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var req *http.Request
			rec := httptest.NewRecorder()
			jsonStr, _ := json.Marshal(tt.request)
			req = httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(jsonStr))
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}

var deleteProductTestCases = []struct {
	tid        int
	name       string
	productID  string
	request    map[string]interface{}
	wantStatus int
	wantBody   string
}{
	{
		tid:        1,
		name:       "商品が正常に削除できること",
		productID:  productIDForTest,
		wantStatus: http.StatusNoContent,
	},
	{
		tid:        2,
		name:       "削除できる商品がない場合は404エラー",
		productID:  "10",
		wantStatus: http.StatusNotFound,
	},
}

func TestDeleteProduct(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE products")
	dbConn.Exec("insert into products (id, product_name, price,remarks,quantity, created_at, updated_at)values('090e142d-baa3-4039-9d21-cf5a1af39094','1',100,'1',1,'2022-06-20 22:14:22','2022-06-20 22:14:22'),('5c3325c1-d539-42d6-b405-2af2f6b99ed9','2',1000,'2',10, '2022-06-20 22:14:23', '2022-06-20 22:14:23');")
	r := gin.New()
	zapLogger, err := logger.NewLoggerForTest(true)
	require.NoError(t, err)
	r.Use(func(_ *gin.Context) { binding.EnableDecoderUseNumber = true })
	r.Use(middleware.NewTracing())
	r.Use(middleware.NewLogging(zapLogger))
	productHandler := NewProductHandler(dbConn, nil, zapLogger)
	r.DELETE("/products/:id", productHandler.DeleteProduct)

	for _, tt := range deleteProductTestCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var req *http.Request
			rec := httptest.NewRecorder()
			jsonStr, _ := json.Marshal(tt.request)
			req = httptest.NewRequest(http.MethodDelete, "/products/"+tt.productID, bytes.NewBuffer(jsonStr))
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, rec.Body.String())
			}
		})
	}
}
