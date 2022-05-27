package auth

import (
	"github.com/AI1411/golang-admin-api/util/jwt"
	"github.com/gin-gonic/gin"
)

func GetLoginUserID(ctx *gin.Context) (string, error) {
	cookie, err := ctx.Cookie("jwt")
	if err != nil {
		return "", err
	}

	id, err := util.ParseJwt(cookie)
	if err != nil {
		return "", err
	}

	return id, nil
}
