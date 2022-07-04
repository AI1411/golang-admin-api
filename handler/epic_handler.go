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

type EpicHandler struct {
	Db     *gorm.DB
	logger *zap.Logger
}

func NewEpicHandler(db *gorm.DB, logger *zap.Logger) *EpicHandler {
	return &EpicHandler{
		Db:     db,
		logger: logger,
	}
}

type epicRequest struct {
	IsOpen          bool   `json:"is_open" binding:"omitempty,boolean" example:"false"`
	AuthorId        string `json:"author_id" binding:"required,uuid4" format:"/^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i" example:"443b5f1c-8a3a-4485-b3bc-05e69b40b290"`
	EpicTitle       string `json:"epic_title" binding:"required,max=64" example:"title"`
	EpicDescription string `json:"epic_description" binding:"omitempty,max=256" example:"description"`
	Label           string `json:"label" binding:"omitempty,max=64" example:"label"`
	MilestoneId     string `json:"milestone_id" binding:"required,uuid4" format:"/^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i" example:"443b5f1c-8a3a-4485-b3bc-05e69b40b290"`
	AssigneeId      string `json:"assignee_id" binding:"omitempty,uuid4" format:"/^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i" example:"443b5f1c-8a3a-4485-b3bc-05e69b40b290"`
	ProjectId       string `json:"project_id" binding:"required,uuid4" format:"/^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i" example:"443b5f1c-8a3a-4485-b3bc-05e69b40b290"`
}

type epicResponseItem struct {
	models.Epic
}

type epicResponse struct {
	Total int                `json:"total"`
	Epics []epicResponseItem `json:"epics"`
}

type searchEpicParams struct {
	IsOpen      string `binding:"omitempty,boolean" form:"is_open"`
	AuthorID    string `binding:"omitempty,len=36" form:"author_id"`
	EpicTitle   string `binding:"omitempty,max=64" form:"epic_title"`
	Label       string `binding:"omitempty,max=64" form:"label"`
	MilestoneID string `binding:"omitempty,len=36" form:"milestone_id"`
	AssigneeID  string `binding:"omitempty,len=36" form:"assignee_id"`
	ProjectID   string `binding:"omitempty,len=36" form:"project_id"`
	Offset      string `form:"offset,default=0" binding:"omitempty,numeric"`
	Limit       string `form:"limit,default=10" binding:"omitempty,numeric"`
}

// GetEpics @title 一覧取得
// @id GetEpics
// @tags epics
// @version バージョン(1.0)
// @description 指定された条件に一致するepic一覧情報を取得する
// @Summary epic一覧取得
// @Produce json
// @Success 200 {object} epicResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /epics [GET]
// @Param is_open query string false "解放フラグ" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param author_id query string false "作成者ID" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param epic_title query string false "タイトル <br><table><tr><th>項目</th><th>説明</th></tr><tr><td>new</td><td>新規</td></tr><tr><td>processing</td><td>進行中</td></tr><tr><td>done</td><td>完了</td></tr><tr><td>closed</td><td>終了</td></tr></table>" Enums(new, processing, done, closed)
// @Param label query string false "ラベル" minlength(36) maxlength(36) format(UUID v4)
// @Param milestone_id query string false "作成日" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param assignee_id query string false "作成日" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param project_id query string false "作成日" format(YYYY-MM-DDThh:mm:ss±hh:mm)
// @Param offset query int false "開始位置" default(0) minimum(0)
// @Param limit query int false "取得上限" default(12) minimum(1) maximum(100)
func (h *EpicHandler) GetEpics(ctx *gin.Context) {
	traceID := appcontext.GetTraceID(ctx)
	var params searchEpicParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	var epics []models.Epic
	query := createEpicQueryBuilder(params, h)
	if err := query.Find(&epics).Error; err != nil {
		h.logger.Error("failed to get epics", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to get epics", err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"total": len(epics),
		"epics": epics,
	})
}

// GetEpicDetail @title epic詳細
// @id GetEpicDetail
// @tags epics
// @version バージョン(1.0)
// @description epic詳細を返す
// @Summary epic詳細取得
// @Produce json
// @Success 200 {object} epicResponseItem
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /epics/:id [GET]
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
func (h *EpicHandler) GetEpicDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	traceID := appcontext.GetTraceID(ctx)
	var epic models.Epic
	if err := h.Db.Where("id = ?", id).First(&epic).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			h.logger.Error("failed to find coupon", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("epic not found"))
		case gorm.ErrInvalidSQL:
			h.logger.Error("invalid sql", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	ctx.JSON(http.StatusOK, epic)
}

// CreateEpic @title epic作成
// @id CreateEpic
// @tags epics
// @version バージョン(1.0)
// @description epicを作成する
// @Summary epic作成
// @Produce json
// @Success 201 {object} epicResponseItem
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /epics [POST]
// @Accept json
// @Param epicRequest body epicRequest true "create epic"
func (h *EpicHandler) CreateEpic(ctx *gin.Context) {
	traceID := appcontext.GetTraceID(ctx)
	var epic models.Epic
	if err := ctx.ShouldBindJSON(&epic); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	if err := h.Db.Create(&epic).Error; err != nil {
		h.logger.Error("failed to create coupon", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to create epic", err))
		return
	}
	ctx.JSON(http.StatusCreated, epic)
}

// UpdateEpic @title epic編集
// @id UpdateEpic
// @tags epics
// @version バージョン(1.0)
// @description epicを編集する
// @Summary epic編集
// @Produce json
// @Success 202 {object} epicResponseItem
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /epics/:id [PUT]
// @Accept json
// @Param epicRequest body epicRequest true "update epic"
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
func (h *EpicHandler) UpdateEpic(ctx *gin.Context) {
	var epic models.Epic
	id := ctx.Param("id")
	traceID := appcontext.GetTraceID(ctx)
	if err := h.Db.Where("id = ?", id).First(&epic).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			h.logger.Error("failed to update epic", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("epic not found"))
		case gorm.ErrInvalidSQL:
			h.logger.Error("invalid sql", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	if err := ctx.ShouldBindJSON(&epic); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	if err := h.Db.Save(&epic).Error; err != nil {
		h.logger.Error("failed to update epic", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to update epic", err))
		return
	}

	ctx.JSON(http.StatusAccepted, epic)
}

// DeleteEpic @title epic削除
// @id DeleteEpic
// @tags epics
// @version バージョン(1.0)
// @description epicを削除する
// @Summary epic削除
// @Produce json
// @Success 204
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /epics/:id [DELETE]
// @Accept json
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
func (h *EpicHandler) DeleteEpic(ctx *gin.Context) {
	var epic models.Epic
	id := ctx.Param("id")
	traceID := appcontext.GetTraceID(ctx)
	if err := h.Db.Where("id = ?", id).First(&epic).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			h.logger.Error("failed to delete epic", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("epic not found"))
		case gorm.ErrInvalidSQL:
			h.logger.Error("invalid sql", zap.Error(err),
				zap.String("trace_id", traceID))
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}

	if err := h.Db.Delete(&epic).Error; err != nil {
		h.logger.Error("failed to delete epic", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to delete epic", err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func createEpicQueryBuilder(params searchEpicParams, h *EpicHandler) *gorm.DB {
	var products []models.Product
	query := h.Db.Find(&products)

	if params.IsOpen != "" {
		query = query.Where("IsOpen = ?", params.IsOpen)
	}
	if params.AuthorID != "" {
		query = query.Where("author_id = ?", params.AuthorID)
	}
	if params.EpicTitle != "" {
		query = query.Where("epic_title LIKE ?", "%"+params.EpicTitle+"%")
	}
	if params.Label != "" {
		query = query.Where("label LIKE ?", "%"+params.Label+"%")
	}
	if params.MilestoneID != "" {
		query = query.Where("milestone_id = ?", params.MilestoneID)
	}
	if params.AssigneeID != "" {
		query = query.Where("assignee_id = ?", params.AssigneeID)
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
