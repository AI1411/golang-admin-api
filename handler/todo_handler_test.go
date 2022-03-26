package handler

import (
	"github.com/AI1411/golang-admin-api/db"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var todoHandlerTestCases = []struct {
	tid        int
	name       string
	request    interface{}
	wantStatus int
	wantBody   string
}{
	{
		tid:        1,
		name:       "TODO一覧が正常に取得できること",
		request:    nil,
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
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/todos", nil))
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}
