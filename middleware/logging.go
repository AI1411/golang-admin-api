package middleware

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/AI1411/golang-admin-api/util/appcontext"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func NewLogging(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := appcontext.GetTraceID(c)
		endpoint := c.Request.RequestURI
		bufBody, _ := io.ReadAll(c.Request.Body)
		reqBody := map[string]interface{}{}
		_ = json.Unmarshal(bufBody, &reqBody)

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bufBody))

		logger.Info("request", zap.String("trace_id", traceID),
			zap.String("http_method", c.Request.Method),
			zap.String("endpoint", endpoint),
			zap.Any("header", c.Request.Header),
			zap.Any("body", reqBody))

		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = writer
		c.Next()

		logger.Info("response", zap.String("trace_id", traceID),
			zap.String("http_method", c.Request.Method),
			zap.String("endpoint", endpoint),
			zap.Any("header", writer.Header()),
			zap.Int("http_status", writer.Status()),
			zap.Any("body", writer.body.String()))
	}
}
