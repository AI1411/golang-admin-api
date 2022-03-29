package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/AI1411/golang-admin-api/models"
	util "github.com/AI1411/golang-admin-api/util/jwt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Claims struct {
	jwt.StandardClaims
}

type AuthHandler struct {
	Db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{Db: db}
}

type authRequest struct {
	LastName             string `json:"last_name" binding:"required"`
	FirstName            string `json:"first_name" binding:"required"`
	Age                  uint8  `json:"age" binding:"required"`
	Image                string `json:"image" binding:"omitempty"`
	Email                string `json:"email" binding:"required"`
	Password             string `json:"password" binding:"required"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var req authRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}

	if req.Password != req.PasswordConfirmation {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "パスワードが一致しません",
		})
		return
	}

	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Image:     req.Image,
		Age:       req.Age,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user.SetPassword(req.Password)
	h.Db.Create(&user)

	ctx.JSON(http.StatusOK, user)
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

	ctx.SetCookie("jwt", token, 3600, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "認証に成功しました",
	})
}

func (h *AuthHandler) Me(ctx *gin.Context) {
	cookie, _ := ctx.Cookie("jwt")

	id, _ := util.ParseJwt(cookie)

	var user models.User

	h.Db.Where("id = ?", id).First(&user)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "its me!",
		"user":    user,
	})
}
