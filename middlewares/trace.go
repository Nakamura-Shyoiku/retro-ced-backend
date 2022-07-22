package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set the trace UUID
		c.Set("trace", uuid.New().String())

		// Continue request
		c.Next()
	}
}
