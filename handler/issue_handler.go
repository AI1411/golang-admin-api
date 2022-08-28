package handler

import (
	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/appcontext"
	"github.com/AI1411/golang-admin-api/util/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"net/http"
)

type IssueHandler struct {
	Db     *gorm.DB
	logger *zap.Logger
}

func NewIssueHandler(db *gorm.DB, logger *zap.Logger) *IssueHandler {
	return &IssueHandler{
		Db:     db,
		logger: logger,
	}
}

type issueRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"omitempty,max=255"`
	UserID      string `json:"user_id" binding:"omitempty"`
	MilestoneID string `json:"milestone_id" binding:"omitempty"`
	IssueStatus string `json:"issue_status" binding:"required"`
}

type issueResponseItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      string `json:"user_id"`
	MilestoneID string `json:"milestone_id"`
	IssueStatus string `json:"issue_status"`
	CreatedAt   string `json:"created_at"`
}

type issueResponse struct {
	Issues []issueResponseItem `json:"items"`
	Total  uint32              `json:"total"`
}

type searchIssueParams struct {
	ID          string `form:"id" binding:"omitempty"`
	Title       string `form:"title" binding:"omitempty,max=64"`
	Description string `form:"description" binding:"omitempty,max=255"`
	UserID      string `form:"assignee_id" binding:"omitempty,len=36" `
	MilestoneID string `form:"milestone_id" binding:"omitempty,len=36" `
	Offset      string `form:"offset,default=0" binding:"omitempty,numeric"`
	Limit       string `form:"limit,default=10" binding:"omitempty,numeric"`
}

// GetIssues @title 一覧取得
// @id GetIssues
// @tags issues
// @version バージョン(1.0)
// @description 指定された条件に一致するissue一覧情報を取得する
// @Summary issue一覧取得
// @Produce json
// @Success 200 {object} issueResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /issues [GET]
// @Param id query string false "ID" minlength(36) maxlength(36) format(UUID v4)
// @Param user_id query string false "ユーザID" minlength(36) maxlength(36) format(UUID v4)
// @Param title query string false "タイトル" minlength(1) maxlength(64)
// @Param description query string false "Issue description" minlength(1) maxlength(255)
// @Param milestone_id query string false "マイルストーンID" minlength(36) maxlength(36) format(UUID v4)
// @Param offset query int false "開始位置" default(0) minimum(0)
// @Param limit query int false "取得上限" default(12) minimum(1) maximum(100)
func (h *IssueHandler) GetIssues(ctx *gin.Context) {
	traceID := appcontext.GetTraceID(ctx)
	var params searchIssueParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	var issues []models.Issue
	query := createIssueQueryBuilder(params, h)
	if err := query.Find(&issues).Error; err != nil {
		h.logger.Error("failed to get issues", zap.Error(err),
			zap.String("trace_id", traceID))
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to get issues", err))
		return
	}
	var response issueResponse
	response.Total = uint32(len(issues))
	if err := copier.Copy(&response.Issues, issues); err != nil {
		h.logger.Error("failed to get issues", zap.Error(err),
			zap.String("trace_id", traceID))
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func createIssueQueryBuilder(params searchIssueParams, h *IssueHandler) *gorm.DB {
	var products []models.Product
	query := h.Db.Find(&products)

	if params.ID != "" {
		query = query.Where("id = ?", params.ID)
	}
	if params.Title != "" {
		query = query.Where("id LIKE ?", "%"+params.Title+"%")
	}
	if params.Description != "" {
		query = query.Where("description LIKE ?", "%"+params.Description+"%")
	}
	if params.UserID != "" {
		query = query.Where("user_id ?", params.ID)
	}
	if params.MilestoneID != "" {
		query = query.Where("milestone_id = ?", params.MilestoneID)
	}
	if params.Offset != "" {
		query = query.Offset(params.Offset)
	}
	if params.Limit != "" {
		query = query.Limit(params.Limit)
	}
	return query
}
