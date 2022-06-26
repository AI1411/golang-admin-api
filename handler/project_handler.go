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
