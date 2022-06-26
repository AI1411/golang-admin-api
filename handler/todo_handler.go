package handler

import (
	"net/http"

	"github.com/AI1411/golang-admin-api/util/appcontext"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
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
	Status    string `form:"status" binding:"omitempty,oneof=success waiting canceled processing done"`
	UserID    string `form:"user_id" binding:"omitempty,max=64"`
	CreatedAt string `form:"created_at" binding:"omitempty,datetime"`
	Offset    string `form:"offset" binding:"omitempty,numeric"`
	Limit     string `form:"limit" binding:"omitempty,numeric"`
}

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
	ctx.JSON(http.StatusOK, gin.H{
		"total": len(todos),
		"todos": todos,
	})
}

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
