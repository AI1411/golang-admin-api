package controllers

import (
	"api/models"
	"api/util/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type TodoHandler struct {
	Db *gorm.DB
}

func (h *TodoHandler) GetAll(ctx *gin.Context) {
	var todos []models.Todo
	h.Db.Preload("User").Find(&todos)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    todos,
	})
}

func (h *TodoHandler) GetDetail(ctx *gin.Context) {
	var todo models.Todo
	id := ctx.Param("id")
	h.Db.Where("id = ?", id).Find(&todo)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    todo,
	})
}

func (h *TodoHandler) CreateTodo(ctx *gin.Context) {
	todo := models.Todo{}
	if err := ctx.ShouldBindJSON(&todo); err != nil {
		restErr := errors.NewBadRequestError("invalid request")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	h.Db.Create(&todo)
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "success",
		"data":    todo,
	})
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
	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "success",
		"data":    todo,
	})
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
	ctx.JSON(http.StatusOK, gin.H{
		"message": "削除されました",
	})
}
