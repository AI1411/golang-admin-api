package handler

import (
	"net/http"

	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type TodoHandler struct {
	Db *gorm.DB
}

func NewTodoHandler(db *gorm.DB) *TodoHandler {
	return &TodoHandler{Db: db}
}

func (h *TodoHandler) GetAll(ctx *gin.Context) {
	var todos []models.Todo
	h.Db.Preload("User").Find(&todos)
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
		restErr := errors.NewBadRequestError(err.Error())
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	h.Db.Create(&todo)
	ctx.JSON(http.StatusCreated, todo)
	return
}

func (h *TodoHandler) UpdateTodo(ctx *gin.Context) {
	todo := models.Todo{}
	id := ctx.Param("id")
	h.Db.First(&todo, id)
	if err := ctx.ShouldBindJSON(&todo); err != nil {
		restErr := errors.NewBadRequestError("invalid request")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	h.Db.Save(&todo)
	ctx.JSON(http.StatusAccepted, todo)
}

func (h TodoHandler) DeleteTodo(ctx *gin.Context) {
	todo := models.Todo{}
	id := ctx.Param("id")
	if err := h.Db.First(&todo, id).Error; err != nil {
		restErr := errors.NewBadRequestError("invalid request")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	h.Db.Delete(&todo)
	ctx.Status(http.StatusNoContent)
}
