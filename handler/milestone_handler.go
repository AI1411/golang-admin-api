package handler

import (
	"net/http"
	"time"

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

type milestoneRequest struct {
	MilestoneTitle       string `json:"milestone_title" example:"milestone title" binding:"required,max=64"`
	MilestoneDescription string `json:"milestone_description" example:"milestone description" binding:"omitempty,max=255"`
	ProjectId            string `json:"project_id" example:"443b5f1c-8a3a-4485-b3bc-05e69b40b290" format:"/^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i" binding:"required,uuid4"`
}

type milestoneResponseItem struct {
	Id                   string    `json:"id" format:"/^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i" example:"milestone id"`
	MilestoneTitle       string    `json:"milestone_title" example:"milestone title"`
	MilestoneDescription string    `json:"milestone_description" example:"milestone description"`
	ProjectId            string    `json:"project_id" format:"/^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i" example:"project id"`
	CreatedAt            time.Time `json:"created_at" example:"2020-01-01T00:00:00Z"`
	UpdatedAt            time.Time `json:"updated_at" example:"2020-01-01T00:00:00Z"`
}

type milestonesResponse struct {
	Total      int                     `json:"total"`
	Milestones []milestoneResponseItem `json:"milestones"`
}

// GetMilestones @title 一覧取得
// @id GetMilestones
// @tags milestones
// @version バージョン(1.0)
// @description 指定された条件に一致するmilestone一覧情報を取得する
// @Summary milestone一覧取得
// @Produce json
// @Success 200 {object} milestonesResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /milestones [GET]
// @Param milestone_title query string false "タイトル" maxlength(64)
// @Param milestone_description query string false "説明" maxlength(255)
// @Param project_id query string false "プロジェクトID" minlength(36) maxlength(36) format(UUID v4)
// @Param offset query int false "開始位置" default(0) minimum(0)
// @Param limit query int false "取得上限" default(12) minimum(1) maximum(100)
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

// GetMilestoneDetail @title milestone詳細
// @id GetMilestoneDetail
// @tags milestones
// @version バージョン(1.0)
// @description milestone詳細を返す
// @Summary milestone詳細取得
// @Produce json
// @Success 200 {object} milestoneResponseItem
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /milestones/:id [GET]
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
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

// CreateMilestone @title milestone作成
// @id CreateMilestone
// @tags milestones
// @version バージョン(1.0)
// @description milestoneを作成する
// @Summary milestone作成
// @Produce json
// @Success 201 {object} milestoneResponseItem
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /milestones [POST]
// @Accept json
// @Param milestoneRequest body milestoneRequest true "create milestone"
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

// UpdateMileStone @title milestone編集
// @id UpdateMileStone
// @tags milestones
// @version バージョン(1.0)
// @description milestoneを編集する
// @Summary milestone編集
// @Produce json
// @Success 202 {object} milestoneResponseItem
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /milestones/:id [PUT]
// @Accept json
// @Param milestoneRequest body milestoneRequest true "update milestone"
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
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

// DeleteMilestone @title milestone削除
// @id DeleteMilestone
// @tags milestones
// @version バージョン(1.0)
// @description milestoneを削除する
// @Summary milestone削除
// @Produce json
// @Success 204
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /milestones/:id [DELETE]
// @Accept json
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
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
