package handler

import (
	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type MilestoneHandler struct {
	Db *gorm.DB
}

func NewMilestoneHandler(db *gorm.DB) *MilestoneHandler {
	return &MilestoneHandler{Db: db}
}

func (h *MilestoneHandler) CreateMilestone(ctx *gin.Context) {
	var milestone models.Milestone
	if err := ctx.ShouldBindJSON(&milestone); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	milestone.CreateUUID()

	if err := h.Db.Create(&milestone).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to create milestone", err))
		return
	}

	ctx.JSON(http.StatusCreated, milestone)
}

func (h *MilestoneHandler) UpdateMileStone(ctx *gin.Context) {
	var milestone models.Milestone
	id := ctx.Param("id")
	if err := h.Db.Where("id = ?", id).First(&milestone).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("milestone not found"))
		case gorm.ErrInvalidSQL:
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	if err := ctx.ShouldBindJSON(&milestone); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	if err := h.Db.Save(&milestone).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to update milestone", err))
		return
	}

	ctx.JSON(http.StatusAccepted, milestone)
}
