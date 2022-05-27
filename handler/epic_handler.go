package handler

import (
	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type EpicHandler struct {
	Db *gorm.DB
}

func NewEpicHandler(db *gorm.DB) *EpicHandler {
	return &EpicHandler{Db: db}
}

func (h *EpicHandler) CreateEpic(ctx *gin.Context) {
	var epic models.Epic
	if err := ctx.ShouldBindJSON(&epic); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	if err := h.Db.Create(&epic).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to create epic", err))
		return
	}
	ctx.JSON(http.StatusCreated, epic)
}

func (h *EpicHandler) UpdateEpic(ctx *gin.Context) {
	var epic models.Epic
	id := ctx.Param("id")
	if err := h.Db.Where("id = ?", id).First(&epic).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("epic not found"))
		case gorm.ErrInvalidSQL:
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	if err := ctx.ShouldBindJSON(&epic); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	if err := h.Db.Save(&epic).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to update epic", err))
		return
	}

	ctx.JSON(http.StatusAccepted, epic)
}
