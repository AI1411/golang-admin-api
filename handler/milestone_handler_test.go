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

const milestoneIDForTest = "de1ccf61-4f17-4f51-8a72-b12d1e5e4191"

var getMilestonesTestCases = []struct {
	tid        int
	name       string
	request    map[string]interface{}
	wantStatus int
	wantBody   string
}{
	{
		tid:        1,
		name:       "マイルストーン一覧が正常に取得できること",
		request:    map[string]interface{}{},
		wantStatus: http.StatusOK,
		wantBody: `{
			"milestones": [
				{
					"id": "de1ccf61-4f17-4f51-8a72-b12d1e5e4191",
					"milestone_title": "なる早",
					"milestone_description": "test description",
					"project_id": "239fabd9-03da-4cd6-bffd-131544f12b5d",
					"created_at": "2022-06-23T07:10:25+09:00",
					"updated_at": "2022-06-23T07:10:25+09:00"
				},
				{
					"id": "86b806e6-f34f-403b-b854-d31a79e2195e",
					"milestone_title": "first",
					"milestone_description": "first",
					"project_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"created_at": "2022-06-23T07:10:47+09:00",
					"updated_at": "2022-06-23T07:10:47+09:00"
				}
			],
			"total": 2
		}`,
	},
	{
		tid:  2,
		name: "検索結果0件",
		request: map[string]interface{}{
			"milestone_title": "test",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"milestones": [],
    		"total": 0
		}`,
	},
	{
		tid:  3,
		name: "パラメータのバリデーションエラー",
		request: map[string]interface{}{
			"milestone_title": strings.Repeat("a", 65),
			"project_id":      "test",
			"offset":          "test",
			"limit":           "test",
		},
		wantStatus: http.StatusBadRequest,
		wantBody: `{
			"code": 400,
			"message": "パラメータが不正です",
			"details": [
				{
					"attribute": "MilestoneTitle",
					"message": "タイトルは不正です"
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
		name: "milestone_title で検索",
		request: map[string]interface{}{
			"milestone_title": "first",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"milestones": [
				{
					"id": "86b806e6-f34f-403b-b854-d31a79e2195e",
					"milestone_title": "first",
					"milestone_description": "first",
					"project_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"created_at": "2022-06-23T07:10:47+09:00",
					"updated_at": "2022-06-23T07:10:47+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  5,
		name: "project_id で範囲検索",
		request: map[string]interface{}{
			"project_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"milestones": [
				{
					"id": "86b806e6-f34f-403b-b854-d31a79e2195e",
					"milestone_title": "first",
					"milestone_description": "first",
					"project_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"created_at": "2022-06-23T07:10:47+09:00",
					"updated_at": "2022-06-23T07:10:47+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  6,
		name: "offset 指定で検索",
		request: map[string]interface{}{
			"offset": 1,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"milestones": [
				{
					"id": "86b806e6-f34f-403b-b854-d31a79e2195e",
					"milestone_title": "first",
					"milestone_description": "first",
					"project_id": "239fabd9-03da-4cd6-bffd-131544f12b5f",
					"created_at": "2022-06-23T07:10:47+09:00",
					"updated_at": "2022-06-23T07:10:47+09:00"
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  7,
		name: "limit 指定で検索",
		request: map[string]interface{}{
			"limit": 1,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"milestones": [
				{
					"id": "de1ccf61-4f17-4f51-8a72-b12d1e5e4191",
					"milestone_title": "なる早",
					"milestone_description": "test description",
					"project_id": "239fabd9-03da-4cd6-bffd-131544f12b5d",
					"created_at": "2022-06-23T07:10:25+09:00",
					"updated_at": "2022-06-23T07:10:25+09:00"
				}
			],
			"total": 1
		}`,
	},
}

func TestGetMilestones(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE milestones")
	require.NoError(t, dbConn.Exec("INSERT INTO milestones (id, milestone_title, milestone_description, project_id, created_at, updated_at) VALUES ('de1ccf61-4f17-4f51-8a72-b12d1e5e4191', 'なる早', 'test description', '239fabd9-03da-4cd6-bffd-131544f12b5d','2022-06-23 07:10:25', '2022-06-23 07:10:25'),('86b806e6-f34f-403b-b854-d31a79e2195e', 'first', 'first', '239fabd9-03da-4cd6-bffd-131544f12b5f','2022-06-23 07:10:47', '2022-06-23 07:10:47');").Error)
	r := gin.New()
	milestoneHandler := NewMilestoneHandler(dbConn)
	r.GET("/milestones", milestoneHandler.GetMilestones)

	for _, tt := range getMilestonesTestCases {
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
			req = httptest.NewRequest(http.MethodGet, "/milestones"+query, nil)
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}

var getMilestoneDetailTestCases = []struct {
	tid         int
	name        string
	milestoneID string
	wantStatus  int
	wantBody    string
}{
	{
		tid:         1,
		name:        "マイルストーン詳細が正常に取得できること",
		milestoneID: milestoneIDForTest,
		wantStatus:  http.StatusOK,
		wantBody: `{
			"id": "de1ccf61-4f17-4f51-8a72-b12d1e5e4191",
			"milestone_title": "なる早",
			"milestone_description": "test description",
			"project_id": "239fabd9-03da-4cd6-bffd-131544f12b5d",
			"created_at": "2022-06-23T07:10:25+09:00",
			"updated_at": "2022-06-23T07:10:25+09:00"
		}`,
	},
	{
		tid:         2,
		name:        "存在しないIDを指定した場合404エラーになること",
		milestoneID: "invalid_milestone",
		wantStatus:  http.StatusNotFound,
		wantBody:    `{"message": "milestone not found","status": 404,"error": "not_found","causes": null}`,
	},
}

func TestMilestoneDetail(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE milestones")
	require.NoError(t, dbConn.Exec("INSERT INTO milestones (id, milestone_title, milestone_description, project_id, created_at, updated_at) VALUES ('de1ccf61-4f17-4f51-8a72-b12d1e5e4191', 'なる早', 'test description', '239fabd9-03da-4cd6-bffd-131544f12b5d','2022-06-23 07:10:25', '2022-06-23 07:10:25'),('86b806e6-f34f-403b-b854-d31a79e2195e', 'first', 'first', '239fabd9-03da-4cd6-bffd-131544f12b5f','2022-06-23 07:10:47', '2022-06-23 07:10:47');").Error)
	r := gin.New()
	milestoneHandler := NewMilestoneHandler(dbConn)
	r.GET("/milestones/:id", milestoneHandler.GetMilestoneDetail)

	for _, tt := range getMilestoneDetailTestCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var req *http.Request
			rec := httptest.NewRecorder()
			req = httptest.NewRequest(http.MethodGet, "/milestones/"+tt.milestoneID, nil)
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}

var deleteMilestoneTestCases = []struct {
	tid         int
	name        string
	milestoneID string
	request     map[string]interface{}
	wantStatus  int
	wantBody    string
}{
	{
		tid:         1,
		name:        "マイルストーンが正常に削除できること",
		milestoneID: milestoneIDForTest,
		wantStatus:  http.StatusNoContent,
	},
	{
		tid:         2,
		name:        "削除できるマイルストーンがない場合は404エラー",
		milestoneID: "10",
		wantStatus:  http.StatusNotFound,
	},
}

func TestDeleteMilestone(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE milestones")
	require.NoError(t, dbConn.Exec("INSERT INTO milestones (id, milestone_title, milestone_description, project_id, created_at, updated_at) VALUES ('de1ccf61-4f17-4f51-8a72-b12d1e5e4191', 'なる早', 'test description', '239fabd9-03da-4cd6-bffd-131544f12b5d','2022-06-23 07:10:25', '2022-06-23 07:10:25'),('86b806e6-f34f-403b-b854-d31a79e2195e', 'first', 'first', '239fabd9-03da-4cd6-bffd-131544f12b5f','2022-06-23 07:10:47', '2022-06-23 07:10:47');").Error)
	r := gin.New()
	milestoneHandler := NewMilestoneHandler(dbConn)
	r.DELETE("/milestones/:id", milestoneHandler.DeleteMilestone)

	for _, tt := range deleteMilestoneTestCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var req *http.Request
			rec := httptest.NewRecorder()
			jsonStr, _ := json.Marshal(tt.request)
			req = httptest.NewRequest(http.MethodDelete, "/milestones/"+tt.milestoneID, bytes.NewBuffer(jsonStr))
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, rec.Body.String())
			}
		})
	}
}
