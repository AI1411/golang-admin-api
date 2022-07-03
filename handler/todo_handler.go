package handler

import (
	"net/http"

	"github.com/AI1411/golang-admin-api/util/appcontext"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"

	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"
)

type TodoHandler struct {
	Db     *gorm.DB
	logger *zap.Logger
}

func NewTodoHandler(db *gorm.DB, logger *zap.Logger) *TodoHandler {
	return &TodoHandler{
		Db:     db,
		logger: logger,
	}
}

type searchTodoPrams struct {
	Title     string `form:"title" binding:"omitempty,max=64"`
	Body      string `form:"body" binding:"omitempty,max=64"`
	Status    string `form:"status" binding:"omitempty,oneof=new processing done closed"`
	UserID    string `form:"user_id" binding:"omitempty,max=64"`
	CreatedAt string `form:"created_at" binding:"omitempty,datetime"`
	Offset    string `form:"offset" binding:"omitempty,numeric"`
	Limit     string `form:"limit" binding:"omitempty,numeric"`
}

type todoRequest struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Status string `json:"status"`
	UserId string `json:"user_id"`
}

type todoItem struct {
	models.Todo
}

type todosResponse struct {
	Total int        `json:"total" example:"1"`
	Todos []todoItem `json:"todos"`
}

// GetAll @title 一覧取得
// @id GetAll
// @tags todos
// @version バージョン(1.0)
// @description 指定された条件に一致するtodo一覧情報を取得する
// @Summary todo一覧取得
// @Produce json
// @Success 200 {object} todosResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /todos [GET]
// @Param title query string false "タイトル" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param body query string false "本文" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param status query string false "ステータス <br><table><tr><th>項目</th><th>説明</th></tr><tr><td>new</td><td>新規</td></tr><tr><td>processing</td><td>進行中</td></tr><tr><td>done</td><td>完了</td></tr><tr><td>closed</td><td>終了</td></tr></table>" Enums(new, processing, done, closed)
// @Param user_id query string false "ユーザID" minlength(36) maxlength(36) format(UUID v4)
// @Param created_at query string false "作成日" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param offset query int false "開始位置" default(0) minimum(0)
// @Param limit query int false "取得上限" default(12) minimum(1) maximum(100)
func (h *TodoHandler) GetAll(ctx *gin.Context) {
	traceID := appcontext.GetTraceID(ctx)
	var params searchTodoPrams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind query params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	var todos []models.Todo
	query := createBaseQueryBuilder(params, h)
	query.Find(&todos)

	res := todosResponse{
		Total: len(todos),
		Todos: []todoItem{},
	}
	res.Total = len(todos)
	if err := copier.Copy(&res.Todos, &todos); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// GetDetail @title todo詳細
// @id GetDetail
// @tags todos
// @version バージョン(1.0)
// @description todo詳細を返す
// @Summary todo詳細取得
// @Produce json
// @Success 200 {object} models.Todo
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /todos/:id [GET]
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
func (h *TodoHandler) GetDetail(ctx *gin.Context) {
	var todo models.Todo
	id := ctx.Param("id")
	if err := h.Db.Where("id = ?", id).Find(&todo).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("todo not found"))
		case gorm.ErrInvalidSQL:
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

// CreateTodo @title todo作成
// @id CreateTodo
// @tags todos
// @version バージョン(1.0)
// @description todoを作成する
// @Summary todo作成
// @Produce json
// @Success 201 {object} todoItem
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /todos [POST]
// @Accept json
// @Param todoRequest body todoRequest true "create todo"
func (h *TodoHandler) CreateTodo(ctx *gin.Context) {
	todo := models.Todo{}
	if err := ctx.ShouldBindJSON(&todo); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	h.Db.Create(&todo)
	ctx.JSON(http.StatusCreated, todo)
}

// UpdateTodo @title todo編集
// @id UpdateTodo
// @tags todos
// @version バージョン(1.0)
// @description todoを編集する
// @Summary todo編集
// @Produce json
// @Success 201 {object} todoItem
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /todos/:id [PUT]
// @Accept json
// @Param todoRequest body todoRequest true "update todo"
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
func (h *TodoHandler) UpdateTodo(ctx *gin.Context) {
	todo := models.Todo{}
	id := ctx.Param("id")
	h.Db.First(&todo, id)
	if err := ctx.ShouldBindJSON(&todo); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	h.Db.Save(&todo)
	ctx.JSON(http.StatusAccepted, todo)
}

// DeleteTodo @title tod削除
// @id DeleteTodo
// @tags todos
// @version バージョン(1.0)
// @description todoを削除する
// @Summary todo編集
// @Produce json
// @Success 204
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /todos/:id [DELETE]
// @Accept json
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
func (h TodoHandler) DeleteTodo(ctx *gin.Context) {
	todo := models.Todo{}
	id := ctx.Param("id")
	if err := h.Db.Where("id = ?", id).First(&todo).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("todo not found"))
		case gorm.ErrInvalidSQL:
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	h.Db.Delete(&todo)
	ctx.Status(http.StatusNoContent)
}

func createBaseQueryBuilder(param searchTodoPrams, h *TodoHandler) *gorm.DB {
	var todos []models.Todo
	query := h.Db.Find(&todos)
	if param.Title != "" {
		query = query.Where("title LIKE ?", "%"+param.Title+"%")
	}
	if param.Body != "" {
		query = query.Where("body LIKE ?", "%"+param.Body+"%")
	}
	if param.Status != "" {
		query = query.Where("status = ?", param.Status)
	}
	if param.UserID != "" {
		query = query.Where("user_id = ?", param.UserID)
	}
	if param.CreatedAt != "" {
		query = query.Where("created_at = ?", param.CreatedAt)
	}
	if param.Offset != "" {
		query = query.Offset(param.Offset)
	}
	if param.Limit != "" {
		query = query.Limit(param.Limit)
	}
	return query
}
