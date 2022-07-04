package handler

import (
	"log"
	"net/http"

	"github.com/AI1411/golang-admin-api/util/appcontext"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"
)

type ProjectHandler struct {
	Db            *gorm.DB
	uuidGenerator models.UUIDGenerator
	logger        *zap.Logger
}

func NewProjectHandler(db *gorm.DB, uuidGenerator models.UUIDGenerator, logger *zap.Logger) *ProjectHandler {
	return &ProjectHandler{
		Db:            db,
		logger:        logger,
		uuidGenerator: uuidGenerator,
	}
}

type searchProjectParams struct {
	ProjectTitle string `form:"project_title" binding:"max=64"`
	Offset       string `form:"offset,default=0" binding:"omitempty,numeric"`
	Limit        string `form:"limit,default=10" binding:"omitempty,numeric"`
}

type projectRequest struct {
	ProjectTitle       string `json:"project_title" binding:"required,max=64" example:"project title"`
	ProjectDescription string `json:"project_description" binding:"required,max=255" example:"project description"`
}

type projectResponseItem struct {
	Id                 string `json:"id" example:"218c51c0-904e-4743-a2ae-94f0e34a0d6f"`
	ProjectTitle       string `json:"project_title" example:"218c51c0-904e-4743-a2ae-94f0e34a0d6f"`
	ProjectDescription string `json:"project_description" example:"description"`
}

type projectsResponse struct {
	Total    int              `json:"total"`
	Projects []models.Project `json:"projects"`
}

// GetProjects @title 一覧取得
// @id GetProjects
// @tags projects
// @version バージョン(1.0)
// @description 指定された条件に一致するproject一覧情報を取得する
// @Summary project一覧取得
// @Produce json
// @Success 200 {object} projectsResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /projects [GET]
// @Param project_title query string false "タイトル" maxlength(64)
// @Param offset query int false "開始位置" default(0) minimum(0)
// @Param limit query int false "取得上限" default(12) minimum(1) maximum(100)
func (h *ProjectHandler) GetProjects(ctx *gin.Context) {
	traceID := appcontext.GetTraceID(ctx)
	var params searchProjectParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	var projects []models.Project
	query := createProjectQueryBuilder(params, h)
	if err := query.Find(&projects).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to get epics", err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"total":    len(projects),
		"projects": projects,
	})
}

// GetProjectDetail @title project詳細
// @id GetProjectDetail
// @tags projects
// @version バージョン(1.0)
// @description project詳細を返す
// @Summary project詳細取得
// @Produce json
// @Success 200 {object} projectResponseItem
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /projects/:id [GET]
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
func (h *ProjectHandler) GetProjectDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	var project models.Project
	if err := h.Db.Where("id = ?", id).Preload("Epics").First(&project).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("project not found"))
		case gorm.ErrInvalidSQL:
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	ctx.JSON(http.StatusOK, project)
}

// CreateProject @title project作成
// @id CreateProject
// @tags projects
// @version バージョン(1.0)
// @description projectを作成する
// @Summary project作成
// @Produce json
// @Success 201 {object} projectResponseItem
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /projects [POST]
// @Accept json
// @Param projectRequest body projectRequest true "create project"
func (h *ProjectHandler) CreateProject(ctx *gin.Context) {
	traceID := appcontext.GetTraceID(ctx)
	var project models.Project
	if err := ctx.ShouldBindJSON(&project); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	project.ID = h.uuidGenerator.GenerateUUID()
	if err := h.Db.Create(&project).Error; err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(201, project)
}

// UpdateProject @title project編集
// @id UpdateProject
// @tags projects
// @version バージョン(1.0)
// @description projectを編集する
// @Summary project編集
// @Produce json
// @Success 202 {object} projectResponseItem
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /projects/:id [PUT]
// @Accept json
// @Param projectRequest body projectRequest true "update project"
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
func (h *ProjectHandler) UpdateProject(ctx *gin.Context) {
	var project models.Project
	id := ctx.Param("id")
	if err := h.Db.Where("id = ?", id).First(&project).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("project not found"))
		case gorm.ErrInvalidSQL:
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	traceID := appcontext.GetTraceID(ctx)
	if err := ctx.ShouldBindJSON(&project); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	if err := h.Db.Save(&project).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to update project", err))
		return
	}

	ctx.JSON(http.StatusAccepted, project)
}

// DeleteProject @title project削除
// @id DeleteProject
// @tags projects
// @version バージョン(1.0)
// @description projectを削除する
// @Summary project削除
// @Produce json
// @Success 204
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /projects/:id [DELETE]
// @Accept json
// @Param id path string true "ID" minlength(36) maxlength(36) format(UUID v4)
func (h *ProjectHandler) DeleteProject(ctx *gin.Context) {
	project := models.Project{}
	id := ctx.Param("id")
	if err := h.Db.Where("id = ?", id).First(&project).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("project not found"))
		case gorm.ErrInvalidSQL:
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	h.Db.Delete(&project)
	ctx.Status(http.StatusNoContent)
}

func createProjectQueryBuilder(param searchProjectParams, h *ProjectHandler) *gorm.DB {
	var projects []models.Project
	query := h.Db.Find(&projects)
	log.Printf("p=%+v", param)
	if param.ProjectTitle != "" {
		query = query.Where("project_title LIKE ?", "%"+param.ProjectTitle+"%")
	}
	if param.Offset != "" {
		query = query.Offset(param.Offset)
	}
	if param.Limit != "" {
		query = query.Limit(param.Limit)
	}
	return query
}
