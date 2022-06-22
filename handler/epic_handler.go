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

func (h *EpicHandler) GetEpics(ctx *gin.Context) {
	var params searchEpicParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}
	var epics []models.Epic
	query := createEpicQueryBuilder(params, h)
	if err := query.Find(&epics).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("failed to get epics", err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"total": len(epics),
		"epics": epics,
	})
}

func (h *EpicHandler) GetEpicDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	var epic models.Epic
	if err := h.Db.Where("id = ?", id).First(&epic).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ctx.JSON(http.StatusNotFound, errors.NewNotFoundError("epic not found"))
		case gorm.ErrInvalidSQL:
			ctx.JSON(http.StatusBadRequest, errors.NewBadRequestError("invalid sql"))
		}
		return
	}
	ctx.JSON(http.StatusOK, epic)
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

func (h *EpicHandler) DeleteEpic(ctx *gin.Context) {
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

	if err := h.Db.Delete(&epic).Error; err != nil {
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
