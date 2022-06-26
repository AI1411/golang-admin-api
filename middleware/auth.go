package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/AI1411/golang-admin-api/db"
	"github.com/AI1411/golang-admin-api/models"
	"github.com/AI1411/golang-admin-api/util/errors"
	"github.com/AI1411/golang-admin-api/util/jwt"
)

func AuthenticateBearer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		bearerToken := strings.Split(token, " ")[1]

		id, err := util.ParseJwt(bearerToken)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized!",
			})
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if id == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized!",
			})
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user models.User
		dbConn := db.Init()

		if err := dbConn.Where("id = ?", id).First(&user).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				ctx.AbortWithStatusJSON(http.StatusNotFound, errors.NewNotFoundError("user not found"))
				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError,
				errors.NewInternalServerError("failed to get user", err))
			return
		}
		ctx.Next()
	}
}
