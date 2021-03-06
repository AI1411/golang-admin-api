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
	if err := h.Db.Transaction(func(tx *gorm.DB) error {
		if err := query.Find(&userGroups).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError,
				errors.NewInternalServerError("failed to get user groups", err))
			return err
		}
		var groupUser []models.GroupUser
		if err := tx.Table("group_user").Find(&groupUser).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError,
				errors.NewInternalServerError("failed to get user groups", err))
			return err
		}
		for i, group := range userGroups {
			var userIDs []string
			for _, user := range groupUser {
				if group.ID == user.GroupID {
					userIDs = append(userIDs, user.UserID)
				}
			}
			var users []models.User
			if err := tx.Table("users").Where("id in (?)", userIDs).Find(&users).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError,
					errors.NewInternalServerError("failed to get user groups", err))
				return err
			}
			if len(users) > 0 {
				group.Users = users
			}
			userGroups[i] = group
		}
		return nil
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError,
			errors.NewInternalServerError("failed to get user groups or group user", err))
		return
	}
	ctx.JSON(http.StatusOK, userGroups)
}

func (h *UserGroupHandler) GetUserGroupsDetail(ctx *gin.Context) {
	var userGroup models.UserGroup
	if err := h.Db.Transaction(func(tx *gorm.DB) error {
		id := ctx.Param("id")
		if err := tx.Where("id = ?", id).First(&userGroup).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError,
				errors.NewInternalServerError("failed to get user groups", err))
			return err
		}
		var groupUser []models.GroupUser
		if err := tx.Table("group_user").Find(&groupUser).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError,
				errors.NewInternalServerError("failed to get user groups", err))
			return err
		}
		var userIDs []string
		for _, user := range groupUser {
			if userGroup.ID == user.GroupID {
				userIDs = append(userIDs, user.UserID)
			}
		}
		var users []models.User
		if err := tx.Table("users").Where("id in (?)", userIDs).Find(&users).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError,
				errors.NewInternalServerError("failed to get user groups", err))
			return err
		}
		if len(users) > 0 {
			userGroup.Users = users
		}
		return nil
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError,
			errors.NewInternalServerError("failed to get user groups or group user", err))
		return
	}
	ctx.JSON(http.StatusOK, userGroup)
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
