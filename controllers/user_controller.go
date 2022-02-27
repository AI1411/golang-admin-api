package controllers

import (
	"api/models"
	"api/util/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type UserHandler struct {
	Db *gorm.DB
}

func (h *UserHandler) GetAllUser(ctx *gin.Context) {
	var users []models.User
	h.Db.Preload("Todos").Find(&users)
	for _, user := range users {
		u := &models.User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Image:     "http://localhost:8084" + user.Image,
			Age:       user.Age,
			Email:     user.Email,
			Password:  user.Password,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Todos:     user.Todos,
		}
		user = *u
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    users,
	})
}

func (h *UserHandler) GetUserDetail(ctx *gin.Context) {
	var user models.User
	id := ctx.Param("id")
	h.Db.Preload("Todos").Where("id = ?", id).Find(&user)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    user,
	})
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	user := models.User{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid request")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	user.SetPassword("123456")

	h.Db.Create(&user)
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "success",
		"data":    user,
	})
	return
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	user := models.User{}
	id := ctx.Param("id")
	h.Db.First(&user, id)
	if err := ctx.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid request")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	h.Db.Save(&user)
	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "success",
		"data":    user,
	})
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	user := models.User{}
	id := ctx.Param("id")
	if err := h.Db.First(&user, id).Error; err != nil {
		restErr := errors.NewBadRequestError("invalid request")
		ctx.JSON(restErr.Status(), restErr)
		return
	}
	h.Db.Delete(&user)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "削除されました",
	})
}
