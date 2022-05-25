package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"
)

type UserGroupHandler struct {
	Db *gorm.DB
}

func NewUserGroupHandler(db *gorm.DB) *UserGroupHandler {
	return &UserGroupHandler{Db: db}
}

type searchUserGroupParams struct {
	GroupName string `form:"group_name" binding:"omitempty,max=64"`
}

func (h *UserGroupHandler) GetAllUserGroups(ctx *gin.Context) {
	var params searchUserGroupParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	var userGroups []models.UserGroup
	query := createUserGroupQueryBuilder(params, h)
	query.Find(&userGroups)
	ctx.JSON(http.StatusOK, userGroups)
}

func (h *UserGroupHandler) CreateUserGroup(ctx *gin.Context) {
	userGroup := models.UserGroup{}
	if err := ctx.ShouldBindJSON(&userGroup); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	userGroup.CreateUUID()
	if err := h.Db.Create(&userGroup).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("Failed to create user group", err))
		return
	}

	ctx.JSON(http.StatusCreated, userGroup)
}

func createUserGroupQueryBuilder(param searchUserGroupParams, h *UserGroupHandler) *gorm.DB {
	var userGroups []models.UserGroup
	query := h.Db.Find(&userGroups)
	if param.GroupName != "" {
		query = query.Where("group_name LIKE ?", "%"+param.GroupName+"%")
	}
	return query
}
