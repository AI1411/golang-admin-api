package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"

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

type createUserGroupParams struct {
	GroupName string   `json:"group_name" binding:"required,max=64"`
	UserIDs   []string `json:"user_ids" binding:"required"`
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
	createUserGroupParams := createUserGroupParams{}
	if err := ctx.ShouldBindJSON(&createUserGroupParams); err != nil {
		res := createValidateErrorResponse(err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	userGroup := models.UserGroup{
		GroupName: createUserGroupParams.GroupName,
	}

	userGroup.CreateUUID()
	if err := h.Db.Transaction(func(tx *gorm.DB) error {
		if err := h.Db.Create(&userGroup).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError,
				errors.NewInternalServerError("Failed to create user group", err))
			return err
		}

		for _, u := range createUserGroupParams.UserIDs {
			groupUser := models.GroupUser{
				UserID:  u,
				GroupID: userGroup.ID,
			}
			if err := h.Db.Table("group_user").Create(&groupUser).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError,
					errors.NewInternalServerError("Failed to create user group user", err))
				return err
			}
		}

		return nil
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError,
			errors.NewInternalServerError("Failed to create user group or group user", err))
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
