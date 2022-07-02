package handler

import (
	"net/http"

	"github.com/AI1411/golang-admin-api/util/appcontext"
	"go.uber.org/zap"

	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type MilestoneHandler struct {
	Db     *gorm.DB
	logger *zap.Logger
}

func NewMilestoneHandler(db *gorm.DB, logger *zap.Logger) *MilestoneHandler {
	return &MilestoneHandler{
		Db:     db,
		logger: logger,
	}
}

type searchMilestoneParams struct {
	MilestoneTitle string `form:"milestone_title" binding:"omitempty,max=64"`
	ProjectID      string `form:"project_id" binding:"omitempty,uuid4" `
	Offset         string `form:"offset,default=0" binding:"omitempty,numeric"`
	Limit          string `form:"limit,default=10" binding:"omitempty,numeric"`
}

func (h *MilestoneHandler) GetMilestones(ctx *gin.Context) {
	traceID := appcontext.GetTraceID(ctx)
	var params searchMilestoneParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	var milestones []models.Milestone
	query := createMilestoneQueryBuilder(params, h)
	if err := query.Find(&milestones).Error; err != nil {
		h.logger.Error("failed to get milestone", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to get milestones", err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"total":      len(milestones),
		"milestones": milestones,
	})
}

func (h *MilestoneHandler) GetMilestoneDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	traceID := appcontext.GetTraceID(ctx)
	var milestone models.Milestone
	if err := h.Db.Where("id = ?", id).First(&milestone).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			h.logger.Error("failed to find milestone", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("milestone not found"))
		case gorm.ErrInvalidSQL:
			h.logger.Error("invalid sql", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	ctx.JSON(http.StatusOK, milestone)
}

func (h *MilestoneHandler) CreateMilestone(ctx *gin.Context) {
	traceID := appcontext.GetTraceID(ctx)
	var milestone models.Milestone
	if err := ctx.ShouldBindJSON(&milestone); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	milestone.CreateUUID()

	if err := h.Db.Create(&milestone).Error; err != nil {
		h.logger.Error("failed to create milestone", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to create milestone", err))
		return
	}

	ctx.JSON(http.StatusCreated, milestone)
}

func (h *MilestoneHandler) UpdateMileStone(ctx *gin.Context) {
	var milestone models.Milestone
	id := ctx.Param("id")
	traceID := appcontext.GetTraceID(ctx)
	if err := h.Db.Where("id = ?", id).First(&milestone).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			h.logger.Error("failed to update milestone", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("milestone not found"))
		case gorm.ErrInvalidSQL:
			h.logger.Error("invalid sql", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	if err := ctx.ShouldBindJSON(&milestone); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	if err := h.Db.Save(&milestone).Error; err != nil {
		h.logger.Error("failed to update milestone", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to update milestone", err))
		return
	}

	ctx.JSON(http.StatusAccepted, milestone)
}

func (h *MilestoneHandler) DeleteMilestone(ctx *gin.Context) {
	var milestone models.Milestone
	id := ctx.Param("id")
	traceID := appcontext.GetTraceID(ctx)
	if err := h.Db.Where("id = ?", id).First(&milestone).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			h.logger.Error("failed to delete milestone", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("milestone not found"))
		case gorm.ErrInvalidSQL:
			h.logger.Error("invalid sql", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}

	if err := h.Db.Delete(&milestone).Error; err != nil {
		h.logger.Error("failed to delete milestone", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to delete milestone", err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func createMilestoneQueryBuilder(params searchMilestoneParams, h *MilestoneHandler) *gorm.DB {
	var products []models.Product
	query := h.Db.Find(&products)

	if params.MilestoneTitle != "" {
		query = query.Where("milestone_title LIKE ?", "%"+params.MilestoneTitle+"%")
	}
	if params.ProjectID != "" {
		query = query.Where("project_id = ?", params.ProjectID)
	}
	if params.Offset != "" {
		query = query.Offset(params.Offset)
	}
	if params.Limit != "" {
		query = query.Limit(params.Limit)
	}
	return query
}
