package handler

import (
	"github.com/AI1411/golang-admin-api/util/appcontext"
	"github.com/AI1411/golang-admin-api/util/redis"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"
	util "github.com/AI1411/golang-admin-api/util/jwt"
)

type Claims struct {
	jwt.StandardClaims
}

type AuthHandler struct {
	Db     *gorm.DB
	logger *zap.Logger
}

func NewAuthHandler(db *gorm.DB, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		Db:     db,
		logger: logger,
	}
}

type authRequest struct {
	LastName             string `json:"last_name" binding:"required"`
	FirstName            string `json:"first_name" binding:"required"`
	Age                  uint8  `json:"age" binding:"required"`
	Email                string `json:"email" binding:"required"`
	Password             string `json:"password" binding:"required"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type meRequest struct {
	JwtToken string `json:"jwt_token" binding:"required"`
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	traceID := appcontext.GetTraceID(ctx)
	var req authRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
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
		Age:       req.Age,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user.CreateUUID()
	user.SetPassword(req.Password)
	if err := h.Db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalServerError("user failed to register", err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	traceID := appcontext.GetTraceID(ctx)
	var req loginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := createValidateErrorResponse(err)
		res.outputErrorLog(h.logger, "failed to bind json params", traceID, err)
		ctx.AbortWithStatusJSON(res.Code, res)
		return
	}

	var user models.User

	err := h.Db.Where("email = ?", req.Email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "ユーザが見つかりませんでした",
		})
		return
	}

	if err = user.ComparePassword(req.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "パスワードが間違っています",
		})
		return
	}

	token, err := util.GenerateJwt(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "認証に失敗しました",
		})
		return
	}
	redis.NewSession(ctx, "jwt", token)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "認証に成功しました",
		"value":   token,
		"user":    user,
	})
}

func (h *AuthHandler) Me(ctx *gin.Context) {
	cookie, err := redis.GetSession(ctx, "jwt")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "認証に失敗しました"})
		return
	}

	id, err := util.ParseJwt(cookie)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized!",
		})
		return
	}

	if id == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized!",
		})
		return
	}

	var user models.User

	if err := h.Db.Where("id = ?", id).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, errors.NewNotFoundError("user not found"))
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			errors.NewInternalServerError("failed to get user", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "its me!",
		"user":    user,
	})
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	redis.DeleteSession(ctx, "jwt")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "logout!",
	})
}
