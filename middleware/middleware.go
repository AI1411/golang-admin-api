package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func Cors() gin.HandlerFunc {
	config := cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowOriginFunc: nil,
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		AllowCredentials: true,
		ExposeHeaders:    nil,
		MaxAge:           24 * time.Hour,
	})

	return config
}
