package controllers

import (
	"api/models"
	util "api/util/jwt"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

type Claims struct {
	jwt.StandardClaims
}

type AuthHandler struct {
	Db *gorm.DB
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var data map[string]string

	if err := ctx.ShouldBindJSON(&data); err != nil {
		panic(err)
	}

	if data["password"] != data["password_confirmation"] {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "パスワードが一致しません",
		})
		return
	}
	age, _ := strconv.Atoi(data["age"])

	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Image:     data["image"],
		Age:       uint8(age),
		Email:     data["email"],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user.SetPassword(data["password"])
	h.Db.Create(&user)

	ctx.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var data map[string]string

	if err := ctx.ShouldBindJSON(&data); err != nil {
		panic(err)
	}

	var user models.User

	err := h.Db.Where("email = ?", data["email"]).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "ユーザが見つかりませんでした",
		})
		return
	}

	if err = user.ComparePassword(data["password"]); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "パスワードが間違っています",
		})
		return
	}

	token, err := util.GenerateJwt(strconv.Itoa(user.ID))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "認証に失敗しました",
		})
		return
	}

	ctx.SetCookie("jwt", token, 3600, "/", "localhost", true, true)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "認証に成功しました",
	})
}
