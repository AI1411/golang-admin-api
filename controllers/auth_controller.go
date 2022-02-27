package controllers

import (
	"api/models"
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
