package handler

import (
	"fmt"
	"github.com/AI1411/golang-admin-api/db"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var todoHandlerTestCases = []struct {
	tid        int
	name       string
	request    map[string]interface{}
	wantStatus int
	wantBody   string
}{
	{
		tid:        1,
		name:       "TODO一覧が正常に取得できること",
		request:    map[string]interface{}{},
		wantStatus: http.StatusOK,
		wantBody: `{
			"todos": [
				{
					"id": 1,
					"title": "test1",
					"body": "body1",
					"status": "success",
					"user_id": 1,
					"created_at": "2022-03-26T21:34:52+09:00",
					"updated_at": "2022-03-26T21:34:52+09:00"
				},
				{
					"id": 2,
					"title": "test2",
					"body": "body2",
					"status": "waiting",
					"user_id": 2,
					"created_at": "2022-03-26T21:34:52+09:00",
					"updated_at": "2022-03-26T21:34:52+09:00"
				}
			],
    		"total": 2
		}`,
	},
	{
		tid:  2,
		name: "検索結果0件",
		request: map[string]interface{}{
			"title": "test3",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"todos": [],
    		"total": 0
		}`,
	},
	{
		tid:  3,
		name: "パラメータのバリデーションエラー",
		request: map[string]interface{}{
			"title":      strings.Repeat("a", 65),
			"body":       strings.Repeat("a", 65),
			"status":     "not_status",
			"user_id":    strings.Repeat("a", 65),
			"created_at": "2020/1/1",
			"offset":     "not_numeric",
			"limit":      "not_numeric",
		},
		wantStatus: http.StatusBadRequest,
		wantBody: `{
			"code": 400,
			"message": "パラメータが不正です",
			"details": [
				{
					"attribute": "Title",
					"message": "Titleは不正です"
				},
				{
					"attribute": "Body",
					"message": "Bodyは不正です"
				},
				{
					"attribute": "Status",
					"message": "Statusは不正です"
				},
 				{
            		"attribute": "UserId",
            		"message": "UserIdは不正です"
				},
				{
					"attribute": "CreatedAt",
					"message": "CreatedAtは不正です"
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
		name: "titleで検索",
		request: map[string]interface{}{
			"title": "test1",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"todos": [
				{
					"id": 1,
					"title": "test1",
					"body": "body1",
					"status": "success",
					"user_id": 1,
					"created_at": "2022-03-26T21:34:52+09:00",
					"updated_at": "2022-03-26T21:34:52+09:00"
				}
			],
    		"total": 1
		}`,
	},
	{
		tid:  5,
		name: "Bodyで検索",
		request: map[string]interface{}{
			"body": "body1",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"todos": [
				{
					"id": 1,
					"title": "test1",
					"body": "body1",
					"status": "success",
					"user_id": 1,
					"created_at": "2022-03-26T21:34:52+09:00",
					"updated_at": "2022-03-26T21:34:52+09:00"
				}
			],
    		"total": 1
		}`,
	},
	{
		tid:  6,
		name: "Statusで検索",
		request: map[string]interface{}{
			"status": "success",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"todos": [
				{
					"id": 1,
					"title": "test1",
					"body": "body1",
					"status": "success",
					"user_id": 1,
					"created_at": "2022-03-26T21:34:52+09:00",
					"updated_at": "2022-03-26T21:34:52+09:00"
				}
			],
    		"total": 1
		}`,
	},
	{
		tid:  7,
		name: "UserIDで検索",
		request: map[string]interface{}{
			"user_id": "1",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"todos": [
				{
					"id": 1,
					"title": "test1",
					"body": "body1",
					"status": "success",
					"user_id": 1,
					"created_at": "2022-03-26T21:34:52+09:00",
					"updated_at": "2022-03-26T21:34:52+09:00"
				}
			],
    		"total": 1
		}`,
	},
}

func TestNewTodoHandler(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE todos")
	dbConn.Exec("insert into todos values (1, 'test1', 'body1', 'success', 1, '2022-03-26 21:34:52', '2022-03-26 21:34:52'),(2, 'test2', 'body2', 'waiting', 2, '2022-03-26 21:34:52', '2022-03-26 21:34:52');")
	r := gin.New()
	todoHandler := NewTodoHandler(dbConn)
	r.GET("/todos", todoHandler.GetAll)

	for _, tt := range todoHandlerTestCases {
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
			req = httptest.NewRequest(http.MethodGet, "/todos"+query, nil)
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}
