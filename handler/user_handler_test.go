package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/AI1411/golang-admin-api/db"
)

var getUsersTestCases = []struct {
	tid        int
	name       string
	request    map[string]interface{}
	wantStatus int
	wantBody   string
}{
	{
		tid:        1,
		name:       "ユーザ一覧が正常に取得できること",
		request:    map[string]interface{}{},
		wantStatus: http.StatusOK,
		wantBody: `{
			"total": 2,
			"users": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"first_name": "1",
					"last_name": "1",
					"age": 22,
					"email": "test@gmail.com",
					"password": "JDJhJDE0JFJDRHc1NGNHY0hNd1cySGJZdFZiOHV0ZUZ1d05jSU5xalBDYmFHM3hMNUszNGhrbmMzdGE2",
					"created_at": "2022-06-20T22:14:22+09:00",
					"updated_at": "2022-06-20T22:14:22+09:00",
					"todos": []
				},
				{
					"id": "5c3325c1-d539-42d6-b405-2af2f6b99ed9",
					"first_name": "2",
					"last_name": "2",
					"age": 37,
					"email": "ishii@gmail.com",
					"password": "JDJhJDE0JFJDRHc1NGNHY0hNd1cySGJZdFZiOHV0ZUZ1d05jSU5xalBDYmFHM3hMNUszNGhrbmMzdGE2",
					"created_at": "2022-06-20T22:14:23+09:00",
					"updated_at": "2022-06-20T22:14:23+09:00",
					"todos": []
				}
			]
		}`,
	},
	{
		tid:  2,
		name: "検索結果0件",
		request: map[string]interface{}{
			"first_name": "failed",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"users": [],
    		"total": 0
		}`,
	},
	{
		tid:  3,
		name: "パラメータのバリデーションエラー",
		request: map[string]interface{}{
			"first_name": strings.Repeat("a", 65),
			"last_name":  strings.Repeat("a", 65),
			"age":        "test",
			"email":      strings.Repeat("a", 65),
			"offset":     "not_numeric",
			"limit":      "not_numeric",
		},
		wantStatus: http.StatusBadRequest,
		wantBody: `{
			"code": 400,
			"message": "パラメータが不正です",
			"details": [
				{
					"attribute": "FirstName",
					"message": "名は不正です"
				},
				{
					"attribute": "LastName",
					"message": "姓は不正です"
				},
				{
					"attribute": "Age",
					"message": "年齢は不正です"
				},
				{
					"attribute": "Email",
					"message": "メールアドレスは不正です"
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
		name: "first nameで検索",
		request: map[string]interface{}{
			"first_name": "1",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"total": 1,
			"users": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"first_name": "1",
					"last_name": "1",
					"age": 22,
					"email": "test@gmail.com",
					"password": "JDJhJDE0JFJDRHc1NGNHY0hNd1cySGJZdFZiOHV0ZUZ1d05jSU5xalBDYmFHM3hMNUszNGhrbmMzdGE2",
					"created_at": "2022-06-20T22:14:22+09:00",
					"updated_at": "2022-06-20T22:14:22+09:00",
					"todos": []
				}
			]
		}`,
	},
	{
		tid:  5,
		name: "last nameで検索",
		request: map[string]interface{}{
			"last_name": "1",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"total": 1,
			"users": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"first_name": "1",
					"last_name": "1",
					"age": 22,
					"email": "test@gmail.com",
					"password": "JDJhJDE0JFJDRHc1NGNHY0hNd1cySGJZdFZiOHV0ZUZ1d05jSU5xalBDYmFHM3hMNUszNGhrbmMzdGE2",
					"created_at": "2022-06-20T22:14:22+09:00",
					"updated_at": "2022-06-20T22:14:22+09:00",
					"todos": []
				}
			]
		}`,
	},
	{
		tid:  6,
		name: "ageで検索",
		request: map[string]interface{}{
			"age": "22",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"total": 1,
			"users": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"first_name": "1",
					"last_name": "1",
					"age": 22,
					"email": "test@gmail.com",
					"password": "JDJhJDE0JFJDRHc1NGNHY0hNd1cySGJZdFZiOHV0ZUZ1d05jSU5xalBDYmFHM3hMNUszNGhrbmMzdGE2",
					"created_at": "2022-06-20T22:14:22+09:00",
					"updated_at": "2022-06-20T22:14:22+09:00",
					"todos": []
				}
			]
		}`,
	},
}

func TestGetUsers(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE users")
	dbConn.Exec("insert into users (id, first_name, last_name, age, email, password, created_at, updated_at)values('090e142d-baa3-4039-9d21-cf5a1af39094','1','1',22,'test@gmail.com','$2a$14$RCDw54cGcHMwW2HbYtVb8uteFuwNcINqjPCbaG3xL5K34hknc3ta6','2022-06-20 22:14:22','2022-06-20 22:14:22'),('5c3325c1-d539-42d6-b405-2af2f6b99ed9','2','2',37,'ishii@gmail.com','$2a$14$RCDw54cGcHMwW2HbYtVb8uteFuwNcINqjPCbaG3xL5K34hknc3ta6', '2022-06-20 22:14:23', '2022-06-20 22:14:23');")
	r := gin.New()
	userHandler := NewUserHandler(dbConn)
	r.GET("/users", userHandler.GetAllUser)

	for _, tt := range getUsersTestCases {
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
			req = httptest.NewRequest(http.MethodGet, "/users"+query, nil)
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}
