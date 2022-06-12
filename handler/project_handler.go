package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"
)

type ProjectHandler struct {
	Db *gorm.DB
}

func NewProjectHandler(db *gorm.DB) *ProjectHandler {
	return &ProjectHandler{Db: db}
}

func (h *ProjectHandler) GetProjects(ctx *gin.Context) {
	var projects []models.Project
	if err := h.Db.Find(&projects).Error; err != nil {
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
	if err := h.Db.Where("id = ?", id).First(&project).Error; err != nil {
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
	var project models.Project
	if err := ctx.ShouldBindJSON(&project); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	project.CreateUUID()
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
	if err := ctx.ShouldBindJSON(&project); err != nil {
		res := createValidateErrorResponse(err)
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
