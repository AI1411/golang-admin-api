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

const projectIDForTest = "090e142d-baa3-4039-9d21-cf5a1af39094"

var getProjectsTestCases = []struct {
	tid        int
	name       string
	request    map[string]interface{}
	wantStatus int
	wantBody   string
}{
	{
		tid:        1,
		name:       "プロジェクト覧が正常に取得できること",
		request:    map[string]interface{}{},
		wantStatus: http.StatusOK,
		wantBody: `{
			"projects": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"project_title": "1",
					"project_description": "1",
					"epics": null
				},
				{
					"id": "5c3325c1-d539-42d6-b405-2af2f6b99ed9",
					"project_title": "2",
					"project_description": "2",
					"epics": null
				}
			],
			"total": 2
		}`,
	},
	{
		tid:  2,
		name: "検索結果0件",
		request: map[string]interface{}{
			"project_title": "failed",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"projects": [],
    		"total": 0
		}`,
	},
	{
		tid:  3,
		name: "パラメータのバリデーションエラー",
		request: map[string]interface{}{
			"project_title": strings.Repeat("a", 65),
			"offset":        "not_numeric",
			"limit":         "not_numeric",
		},
		wantStatus: http.StatusBadRequest,
		wantBody: `{
			"code": 400,
			"message": "パラメータが不正です",
			"details": [
				{
					"attribute": "ProjectTitle",
					"message": "プロジェクト名は不正です"
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
		name: "project_title で検索",
		request: map[string]interface{}{
			"project_title": "1",
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"projects": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"project_title": "1",
					"project_description": "1",
					"epics": null
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  5,
		name: "offset指定で検索",
		request: map[string]interface{}{
			"offset": 1,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"projects": [
				{
					"id": "5c3325c1-d539-42d6-b405-2af2f6b99ed9",
					"project_title": "2",
					"project_description": "2",
					"epics": null
				}
			],
			"total": 1
		}`,
	},
	{
		tid:  6,
		name: "limit指定で検索",
		request: map[string]interface{}{
			"limit": 1,
		},
		wantStatus: http.StatusOK,
		wantBody: `{
			"projects": [
				{
					"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
					"project_title": "1",
					"project_description": "1",
					"epics": null
				}
			],
			"total": 1
		}`,
	},
}

func TestGetProjects(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE projects")
	dbConn.Exec("insert into projects (id, project_title, project_description, created_at, updated_at)values('090e142d-baa3-4039-9d21-cf5a1af39094','1','1','2022-06-20 22:14:22','2022-06-20 22:14:22'),('5c3325c1-d539-42d6-b405-2af2f6b99ed9','2','2', '2022-06-20 22:14:23', '2022-06-20 22:14:23');")
	r := gin.New()
	zapLogger, err := logger.NewLoggerForTest(true)
	require.NoError(t, err)
	r.Use(func(_ *gin.Context) { binding.EnableDecoderUseNumber = true })
	r.Use(middleware.NewTracing())
	r.Use(middleware.NewLogging(zapLogger))
	projectHandler := NewProjectHandler(dbConn, nil, zapLogger)
	r.GET("/projects", projectHandler.GetProjects)

	for _, tt := range getProjectsTestCases {
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
			req = httptest.NewRequest(http.MethodGet, "/projects"+query, nil)
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}

var getProjectDetailTestCases = []struct {
	tid        int
	name       string
	projectID  string
	wantStatus int
	wantBody   string
}{
	{
		tid:        1,
		name:       "プロジェクト詳細が正常に取得できること",
		projectID:  projectIDForTest,
		wantStatus: http.StatusOK,
		wantBody: `{
			"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
			"project_title": "1",
			"project_description": "1",
			"epics": []
		}`,
	},
	{
		tid:        2,
		name:       "存在しないIDを指定した場合404エラーになること",
		projectID:  "invalid_project",
		wantStatus: http.StatusNotFound,
		wantBody:   `{"message": "project not found","status": 404,"error": "not_found","causes": null}`,
	},
}

func TestProjectDetail(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE projects")
	dbConn.Exec("insert into projects (id, project_title, project_description, created_at, updated_at)values('090e142d-baa3-4039-9d21-cf5a1af39094','1','1','2022-06-20 22:14:22','2022-06-20 22:14:22'),('5c3325c1-d539-42d6-b405-2af2f6b99ed9','2','2', '2022-06-20 22:14:23', '2022-06-20 22:14:23');")
	r := gin.New()
	zapLogger, err := logger.NewLoggerForTest(true)
	require.NoError(t, err)
	r.Use(func(_ *gin.Context) { binding.EnableDecoderUseNumber = true })
	r.Use(middleware.NewTracing())
	r.Use(middleware.NewLogging(zapLogger))
	projectHandler := NewProjectHandler(dbConn, nil, zapLogger)
	r.GET("/projects/:id", projectHandler.GetProjectDetail)

	for _, tt := range getProjectDetailTestCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var req *http.Request
			rec := httptest.NewRecorder()
			req = httptest.NewRequest(http.MethodGet, "/projects/"+tt.projectID, nil)
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}

var createProjectTestCases = []struct {
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
			"project_title":       "project title",
			"project_description": "test",
		},
		wantStatus: http.StatusCreated,
		wantBody: `{
			"id": "090e142d-baa3-4039-9d21-cf5a1af39094",
			"project_title": "project title",
			"project_description": "test",
			"epics": null
		}`,
	},
	{
		tid:  2,
		name: "バリデーションエラー",
		request: map[string]interface{}{
			"project_title":       strings.Repeat("a", 65),
			"project_description": strings.Repeat("a", 256),
		},
		wantStatus: http.StatusBadRequest,
		wantBody: `{
			"code": 400,
			"message": "パラメータが不正です",
			"details": [
				{
					"attribute": "ProjectTitle",
					"message": "プロジェクト名は不正です"
				},
				{
					"attribute": "ProjectDescription",
					"message": "ProjectDescriptionは不正です"
				}
			]
		}`,
	},
}

func TestCreateProject(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE projects")
	r := gin.New()
	zapLogger, err := logger.NewLoggerForTest(true)
	require.NoError(t, err)
	r.Use(func(_ *gin.Context) { binding.EnableDecoderUseNumber = true })
	r.Use(middleware.NewTracing())
	r.Use(middleware.NewLogging(zapLogger))
	mockCtrl := gomock.NewController(t)
	uuidGen := mock_models.NewMockUUIDGenerator(mockCtrl)
	uuidGen.EXPECT().GenerateUUID().Return(projectIDForTest)

	projectHandler := NewProjectHandler(dbConn, uuidGen, zapLogger)
	r.POST("/projects", projectHandler.CreateProject)

	for _, tt := range createProjectTestCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var req *http.Request
			rec := httptest.NewRecorder()
			jsonStr, _ := json.Marshal(tt.request)
			req = httptest.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer(jsonStr))
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}

var deleteProjectTestCases = []struct {
	tid        int
	name       string
	projectID  string
	request    map[string]interface{}
	wantStatus int
	wantBody   string
}{
	{
		tid:        1,
		name:       "プロジェクトが正常に削除できること",
		projectID:  projectIDForTest,
		wantStatus: http.StatusNoContent,
	},
	{
		tid:        2,
		name:       "削除できるプロジェクトがない場合は404エラー",
		projectID:  "10",
		wantStatus: http.StatusNotFound,
	},
}

func TestDeleteProject(t *testing.T) {
	dbConn := db.Init()
	dbConn.Exec("TRUNCATE TABLE projects")
	dbConn.Exec("insert into projects (id, project_title, project_description, created_at, updated_at)values('090e142d-baa3-4039-9d21-cf5a1af39094','1','1','2022-06-20 22:14:22','2022-06-20 22:14:22'),('5c3325c1-d539-42d6-b405-2af2f6b99ed9','2','2', '2022-06-20 22:14:23', '2022-06-20 22:14:23');")
	r := gin.New()
	zapLogger, err := logger.NewLoggerForTest(true)
	require.NoError(t, err)
	r.Use(func(_ *gin.Context) { binding.EnableDecoderUseNumber = true })
	r.Use(middleware.NewTracing())
	r.Use(middleware.NewLogging(zapLogger))
	projectHandler := NewProjectHandler(dbConn, nil, zapLogger)
	r.DELETE("/projects/:id", projectHandler.DeleteProject)

	for _, tt := range deleteProjectTestCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var req *http.Request
			rec := httptest.NewRecorder()
			jsonStr, _ := json.Marshal(tt.request)
			req = httptest.NewRequest(http.MethodDelete, "/projects/"+tt.projectID, bytes.NewBuffer(jsonStr))
			r.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, rec.Body.String())
			}
		})
	}
}
