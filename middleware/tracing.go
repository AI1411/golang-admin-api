package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/AI1411/golang-admin-api/util/appcontext"
)

func NewTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := findTraceIDFromHeader(c)
		if traceID == "" {
			var err error
			traceID, err = appcontext.GenerateTraceID()
			if err != nil {
				traceID = "failed to generate traceID"
			}
		}

		appcontext.SetTraceIDIntoContext(c, traceID)

		c.Next()
	}
}

func findTraceIDFromHeader(c *gin.Context) string {
	keys := []string{"x-cgi-trace-id", "x-trace-id"}

	traceID := ""
	for _, key := range keys {
		traceID = c.GetHeader(key)
		if traceID != "" {
			break
		}
	}

	return traceID
}
