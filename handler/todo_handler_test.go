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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/AI1411/golang-admin-api/db"
	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util"
)

var getTodosTestCases = []struct {
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
					"message": "タイトルは不正です"
				},
				{
					"attribute": "Body",
					"message": "本文は不正です"
				},
				{
					"attribute": "Status",
					"message": "ステータスは不正です"
				},
 				{
            		"attribute": "UserID",
            		"message": "ユーザーIDは不正です"
				},
				{
					"attribute": "CreatedAt",
					"message": "作成日時は不正です"
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

func TestGetTodos(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE todos")
	dbConn.Exec("insert into todos values (1, 'test1', 'body1', 'success', 1, '2022-03-26 21:34:52', '2022-03-26 21:34:52'),(2, 'test2', 'body2', 'waiting', 2, '2022-03-26 21:34:52', '2022-03-26 21:34:52');")
	r := gin.New()
	todoHandler := NewTodoHandler(dbConn)
	r.GET("/todos", todoHandler.GetAll)

	for _, tt := range getTodosTestCases {
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

var getTodoDetailTestCases = []struct {
	tid        int
	name       string
	todoID     string
	wantStatus int
	wantBody   string
}{
	{
		tid:        1,
		name:       "TODO詳細が正常に取得できること",
		todoID:     "1",
		wantStatus: http.StatusOK,
		wantBody: `{
				"id": 1,
				"title": "test1",
				"body": "body1",
				"status": "success",
				"user_id": 1,
				"created_at": "2022-03-26T21:34:52+09:00",
				"updated_at": "2022-03-26T21:34:52+09:00"
		}`,
	},
	{
		tid:        2,
		name:       "存在しないIDを指定した場合404エラーになること",
		todoID:     "10",
		wantStatus: http.StatusNotFound,
		wantBody:   `{"message": "todo not found","status": 404,"error": "not_found","causes": null}`,
	},
}

func TestTodoDetail(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE todos")
	dbConn.Exec("insert into todos values (1, 'test1', 'body1', 'success', 1, '2022-03-26 21:34:52', '2022-03-26 21:34:52'),(2, 'test2', 'body2', 'waiting', 2, '2022-03-26 21:34:52', '2022-03-26 21:34:52');")
	r := gin.New()
	todoHandler := NewTodoHandler(dbConn)
	r.GET("/todos/:id", todoHandler.GetDetail)

	for _, tt := range getTodoDetailTestCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var req *http.Request
			rec := httptest.NewRecorder()
			req = httptest.NewRequest(http.MethodGet, "/todos/"+tt.todoID, nil)
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}

var createTodoTestCases = []struct {
	tid        int
	name       string
	request    map[string]interface{}
	wantStatus int
	wantBody   string
}{
	{
		tid:  1,
		name: "TODOが正常に作成できること",
		request: map[string]interface{}{
			"title":      "test",
			"body":       "test",
			"status":     "success",
			"user_id":    1,
			"created_at": "2022-01-01T00:00:00.880012+09:00",
			"updated_at": "2022-01-01T00:00:00.880012+09:00",
		},
		wantStatus: http.StatusCreated,
		wantBody: `{
			"id": 1,
			"title": "test",
			"body": "test",
			"status": "success",
			"user_id": 1,
			"created_at": "2022-01-01T00:00:00.880012+09:00",
			"updated_at": "2022-01-01T00:00:00.880012+09:00"
		}`,
	},
	{
		tid:  2,
		name: "バリデーションエラー",
		request: map[string]interface{}{
			"title":  strings.Repeat("a", 65),
			"body":   strings.Repeat("a", 65),
			"status": "invalid_status",
		},
		wantStatus: http.StatusBadRequest,
		wantBody: `{
			"code": 400,
			"message": "パラメータが不正です",
			"details": [
				{
					"attribute": "Title",
					"message": "タイトルは不正です"
				},
				{
					"attribute": "Body",
					"message": "本文は不正です"
				},
				{
					"attribute": "Status",
					"message": "ステータスは不正です"
				},
 				{
            		"attribute": "UserID",
            		"message": "ユーザーIDは必須です"
				}
			]
		}`,
	},
}

func TestCreateTodo(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE todos")
	r := gin.New()
	todoHandler := NewTodoHandler(dbConn)
	r.POST("/todos", todoHandler.CreateTodo)

	for _, tt := range createTodoTestCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var req *http.Request
			rec := httptest.NewRecorder()
			jsonStr, _ := json.Marshal(tt.request)
			req = httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(jsonStr))
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}

var updateTodoTestCases = []struct {
	tid            int
	name           string
	todoID         string
	request        map[string]interface{}
	wantStatus     int
	wantBody       string
	checkUpdatedAt bool
}{
	{
		tid:    1,
		name:   "TODOが正常に更新できること",
		todoID: "1",
		request: map[string]interface{}{
			"title":   "updated",
			"body":    "updated",
			"status":  "waiting",
			"user_id": 1,
		},
		wantStatus:     http.StatusAccepted,
		checkUpdatedAt: true,
	},
	{
		tid:    2,
		name:   "バリデーションエラー",
		todoID: "1",
		request: map[string]interface{}{
			"title":   strings.Repeat("a", 65),
			"body":    strings.Repeat("a", 65),
			"status":  "invalid_status",
			"user_id": nil,
		},
		wantStatus: http.StatusBadRequest,
		wantBody: `{
			"code": 400,
			"message": "パラメータが不正です",
			"details": [
				{
					"attribute": "Title",
					"message": "タイトルは不正です"
				},
				{
					"attribute": "Body",
					"message": "本文は不正です"
				},
				{
					"attribute": "Status",
					"message": "ステータスは不正です"
				},
 				{
            		"attribute": "UserID",
            		"message": "ユーザーIDは必須です"
				}
			]
		}`,
	},
}

func TestUpdateTodo(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE todos")
	dbConn.Exec("insert into todos values (1, 'test1', 'body1', 'success', 1, '2022-03-26 21:34:52', '2022-03-26 21:34:52'),(2, 'test2', 'body2', 'waiting', 2, '2022-03-26 21:34:52', '2022-03-26 21:34:52');")
	r := gin.New()
	todoHandler := NewTodoHandler(dbConn)
	r.PUT("/todos/:id", todoHandler.UpdateTodo)

	for _, tt := range updateTodoTestCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var req *http.Request
			rec := httptest.NewRecorder()
			jsonStr, _ := json.Marshal(tt.request)
			req = httptest.NewRequest(http.MethodPut, "/todos/"+tt.todoID, bytes.NewBuffer(jsonStr))
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, rec.Body.String())
			}
			if tt.checkUpdatedAt {
				var todo models.Todo
				dbConn.Where("id = 1").First(&todo)
				assert.Equal(t, "updated", todo.Title)
				assert.Equal(t, "updated", todo.Body)
				assert.Equal(t, "waiting", todo.Status)
				assert.Equal(t, util.Uint64ToPtr(1), todo.UserID)
			}
		})
	}
}

var deleteTodoTestCases = []struct {
	tid        int
	name       string
	todoID     string
	request    map[string]interface{}
	wantStatus int
	wantBody   string
}{
	{
		tid:        1,
		name:       "TODOが正常に削除できること",
		todoID:     "1",
		wantStatus: http.StatusNoContent,
	},
	{
		tid:        2,
		name:       "削除できるTODOが場合は404エラー",
		todoID:     "10",
		wantStatus: http.StatusNotFound,
	},
}

func TestDeleteTodo(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE todos")
	dbConn.Exec("insert into todos values (1, 'test1', 'body1', 'success', 1, '2022-03-26 21:34:52', '2022-03-26 21:34:52'),(2, 'test2', 'body2', 'waiting', 2, '2022-03-26 21:34:52', '2022-03-26 21:34:52');")
	r := gin.New()
	todoHandler := NewTodoHandler(dbConn)
	r.DELETE("/todos/:id", todoHandler.DeleteTodo)

	for _, tt := range deleteTodoTestCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var req *http.Request
			rec := httptest.NewRecorder()
			jsonStr, _ := json.Marshal(tt.request)
			req = httptest.NewRequest(http.MethodDelete, "/todos/"+tt.todoID, bytes.NewBuffer(jsonStr))
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, rec.Body.String())
			}
		})
	}
}
