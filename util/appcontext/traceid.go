package appcontext

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const traceIDKey = "api-trace-id"

func SetTraceIDIntoContext(c *gin.Context, traceID string) {
	c.Set(traceIDKey, traceID)
}

func GetTraceID(c *gin.Context) string {
	traceID, exists := c.Get(traceIDKey)
	if !exists {
		traceID, err := GenerateTraceID()
		if err != nil {
			return "failed to generate traceID"
		}
		return traceID
	}

	return traceID.(string)
}

func GenerateTraceID() (string, error) {
	traceUUID, err := uuid.NewRandom()
	traceID := traceUUID.String()
	if err != nil {
		return traceID, err
	}

	return traceID, nil
}
