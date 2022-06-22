package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/AI1411/golang-admin-api/db"
)

var getEpicsTestCases = []struct {
	tid        int
	name       string
	request    map[string]interface{}
	wantStatus int
	wantBody   string
}{
	{
		tid:        1,
		name:       "エピック一覧が正常に取得できること",
		request:    map[string]interface{}{},
		wantStatus: http.StatusOK,
		wantBody: `{
			"epics": [
				{
					"id": 1,
					"is_open": true,
					"author_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"epic_title": "epic",
					"epic_description": "epic",
					"label": "epic",
					"milestone_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"assignee_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"project_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"created_at": "2022-06-22T07:48:17+09:00",
					"updated_at": "2022-06-22T07:48:17+09:00"
				},
				{
					"id": 2,
					"is_open": true,
					"author_id": "239fabd9-03da-4cd6-bffd-131544f12b5d",
					"epic_title": "test",
					"epic_description": "test",
					"label": "test",
					"milestone_id": "239fabd9-03da-4cd6-bffd-131544f12b5d",
					"assignee_id": "239fabd9-03da-4cd6-bffd-131544f12b5d",
					"project_id": "239fabd9-03da-4cd6-bffd-131544f12b5d",
					"created_at": "2022-06-22T07:48:19+09:00",
					"updated_at": "2022-06-22T07:48:19+09:00"
				}
			],
			"total": 2
		}`,
	},
	{
		tid:  2,
		name: "検索結果0件",
		request: map[string]interface{}{
			"label": userIDForTest,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"epics": [],
    		"total": 0
		}`,
	},
	{
		tid:  3,
		name: "パラメータのバリデーションエラー",
		request: map[string]interface{}{
			"is_open":      "test",
			"author_id":    "test",
			"epic_title":   "test",
			"label":        strings.Repeat("a", 65),
			"milestone_id": "test",
			"assignee_id":  "test",
			"project_id":   "test",
			"offset":       "test",
			"limit":        "test",
		},
		wantStatus: http.StatusBadRequest,
		wantBody: `{
			"code": 400,
			"message": "パラメータが不正です",
			"details": [
				{
					"attribute": "IsOpen",
					"message": "IsOpenは不正です"
				},
				{
					"attribute": "AuthorID",
					"message": "AuthorIDは不正です"
				},
				{
					"attribute": "Label",
					"message": "Labelは不正です"
				},
				{
					"attribute": "MilestoneID",
					"message": "MilestoneIDは不正です"
				},
				{
					"attribute": "AssigneeID",
					"message": "AssigneeIDは不正です"
				},
				{
					"attribute": "ProjectID",
					"message": "プロジェクトIDは不正です"
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
		name: "label で検索",
		request: map[string]interface{}{
			"label": "test",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
    		"epics": [
				{
					"id": 2,
					"is_open": true,
					"author_id": "239fabd9-03da-4cd6-bffd-131544f12b5d",
					"epic_title": "test",
					"epic_description": "test",
					"label": "test",
					"milestone_id": "239fabd9-03da-4cd6-bffd-131544f12b5d",
					"assignee_id": "239fabd9-03da-4cd6-bffd-131544f12b5d",
					"project_id": "239fabd9-03da-4cd6-bffd-131544f12b5d",
					"created_at": "2022-06-22T07:48:19+09:00",
					"updated_at": "2022-06-22T07:48:19+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  5,
		name: "milestone_id で範囲検索",
		request: map[string]interface{}{
			"milestone_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"epics": [
				{
					"id": 1,
					"is_open": true,
					"author_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"epic_title": "epic",
					"epic_description": "epic",
					"label": "epic",
					"milestone_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"assignee_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"project_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"created_at": "2022-06-22T07:48:17+09:00",
					"updated_at": "2022-06-22T07:48:17+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  6,
		name: "assignee_id で検索",
		request: map[string]interface{}{
			"assignee_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"epics": [
				{
					"id": 1,
					"is_open": true,
					"author_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"epic_title": "epic",
					"epic_description": "epic",
					"label": "epic",
					"milestone_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"assignee_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"project_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"created_at": "2022-06-22T07:48:17+09:00",
					"updated_at": "2022-06-22T07:48:17+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  7,
		name: "project_id で検索",
		request: map[string]interface{}{
			"project_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"epics": [
				{
					"id": 1,
					"is_open": true,
					"author_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"epic_title": "epic",
					"epic_description": "epic",
					"label": "epic",
					"milestone_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"assignee_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"project_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"created_at": "2022-06-22T07:48:17+09:00",
					"updated_at": "2022-06-22T07:48:17+09:00"
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
			"epics": [
				{
					"id": 2,
					"is_open": true,
					"author_id": "239fabd9-03da-4cd6-bffd-131544f12b5d",
					"epic_title": "test",
					"epic_description": "test",
					"label": "test",
					"milestone_id": "239fabd9-03da-4cd6-bffd-131544f12b5d",
					"assignee_id": "239fabd9-03da-4cd6-bffd-131544f12b5d",
					"project_id": "239fabd9-03da-4cd6-bffd-131544f12b5d",
					"created_at": "2022-06-22T07:48:19+09:00",
					"updated_at": "2022-06-22T07:48:19+09:00"
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
			"epics": [
				{
					"id": 1,
					"is_open": true,
					"author_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"epic_title": "epic",
					"epic_description": "epic",
					"label": "epic",
					"milestone_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"assignee_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"project_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"created_at": "2022-06-22T07:48:17+09:00",
					"updated_at": "2022-06-22T07:48:17+09:00"
				}
			],
			"total": 1
		}`,
	},
}

func TestGetEpics(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE epics")
	require.NoError(t, dbConn.Exec("INSERT INTO go.epics (id, is_open, author_id, epic_title, epic_description, label, milestone_id, assignee_id,project_id, created_at, updated_at) VALUES (1, 1, '239fabd9-03da-4cd6-bffd-131544f12b5f', 'epic', 'epic', 'epic', '239fabd9-03da-4cd6-bffd-131544f12b5f','239fabd9-03da-4cd6-bffd-131544f12b5f', '239fabd9-03da-4cd6-bffd-131544f12b5f', '2022-06-22 07:48:17','2022-06-22 07:48:17'),(2, 1, '239fabd9-03da-4cd6-bffd-131544f12b5d', 'test', 'test', 'test', '239fabd9-03da-4cd6-bffd-131544f12b5d','239fabd9-03da-4cd6-bffd-131544f12b5d', '239fabd9-03da-4cd6-bffd-131544f12b5d', '2022-06-22 07:48:19','2022-06-22 07:48:19');").Error)
	r := gin.New()
	epicHandler := NewEpicHandler(dbConn)
	r.GET("/epics", epicHandler.GetEpics)

	for _, tt := range getEpicsTestCases {
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
			req = httptest.NewRequest(http.MethodGet, "/epics"+query, nil)
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}

var getEpicDetailTestCases = []struct {
	tid        int
	name       string
	epicID     string
	wantStatus int
	wantBody   string
}{
	{
		tid:        1,
		name:       "エピック詳細が正常に取得できること",
		epicID:     "1",
		wantStatus: http.StatusOK,
		wantBody: `{
			"id": 1,
			"is_open": true,
			"author_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
			"epic_title": "epic",
			"epic_description": "epic",
			"label": "epic",
			"milestone_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
			"assignee_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
			"project_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
			"created_at": "2022-06-22T07:48:17+09:00",
			"updated_at": "2022-06-22T07:48:17+09:00"
		}`,
	},
	{
		tid:        2,
		name:       "存在しないIDを指定した場合404エラーになること",
		epicID:     "invalid_epic",
		wantStatus: http.StatusNotFound,
		wantBody:   `{"message": "epic not found","status": 404,"error": "not_found","causes": null}`,
	},
}

func TestEpicDetail(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE epics")
	require.NoError(t, dbConn.Exec("INSERT INTO go.epics (id, is_open, author_id, epic_title, epic_description, label, milestone_id, assignee_id,project_id, created_at, updated_at) VALUES (1, 1, '239fabd9-03da-4cd6-bffd-131544f12b5f', 'epic', 'epic', 'epic', '239fabd9-03da-4cd6-bffd-131544f12b5f','239fabd9-03da-4cd6-bffd-131544f12b5f', '239fabd9-03da-4cd6-bffd-131544f12b5f', '2022-06-22 07:48:17','2022-06-22 07:48:17'),(2, 1, '239fabd9-03da-4cd6-bffd-131544f12b5d', 'test', 'test', 'test', '239fabd9-03da-4cd6-bffd-131544f12b5d','239fabd9-03da-4cd6-bffd-131544f12b5d', '239fabd9-03da-4cd6-bffd-131544f12b5d', '2022-06-22 07:48:19','2022-06-22 07:48:19');").Error)
	r := gin.New()
	epicHandler := NewEpicHandler(dbConn)
	r.GET("/epics/:id", epicHandler.GetEpicDetail)

	for _, tt := range getEpicDetailTestCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var req *http.Request
			rec := httptest.NewRecorder()
			req = httptest.NewRequest(http.MethodGet, "/epics/"+tt.epicID, nil)
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}

var deleteEpicTestCases = []struct {
	tid        int
	name       string
	epicID     string
	request    map[string]interface{}
	wantStatus int
	wantBody   string
}{
	{
		tid:        1,
		name:       "エピックが正常に削除できること",
		epicID:     "1",
		wantStatus: http.StatusNoContent,
	},
	{
		tid:        2,
		name:       "削除できるエピックがない場合は404エラー",
		epicID:     "10",
		wantStatus: http.StatusNotFound,
	},
}

func TestDeleteEpic(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE epics")
	require.NoError(t, dbConn.Exec("INSERT INTO go.epics (id, is_open, author_id, epic_title, epic_description, label, milestone_id, assignee_id,project_id, created_at, updated_at) VALUES (1, 1, '239fabd9-03da-4cd6-bffd-131544f12b5f', 'epic', 'epic', 'epic', '239fabd9-03da-4cd6-bffd-131544f12b5f','239fabd9-03da-4cd6-bffd-131544f12b5f', '239fabd9-03da-4cd6-bffd-131544f12b5f', '2022-06-22 07:48:17','2022-06-22 07:48:17'),(2, 1, '239fabd9-03da-4cd6-bffd-131544f12b5d', 'test', 'test', 'test', '239fabd9-03da-4cd6-bffd-131544f12b5d','239fabd9-03da-4cd6-bffd-131544f12b5d', '239fabd9-03da-4cd6-bffd-131544f12b5d', '2022-06-22 07:48:19','2022-06-22 07:48:19');").Error)
	r := gin.New()
	epicHandler := NewEpicHandler(dbConn)
	r.DELETE("/epics/:id", epicHandler.DeleteEpic)

	for _, tt := range deleteEpicTestCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var req *http.Request
			rec := httptest.NewRecorder()
			jsonStr, _ := json.Marshal(tt.request)
			req = httptest.NewRequest(http.MethodDelete, "/epics/"+tt.epicID, bytes.NewBuffer(jsonStr))
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, rec.Body.String())
			}
		})
	}
}
