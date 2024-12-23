package middleware

import (
	"context"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type ContextKey string

const RequestIDKey ContextKey = "RequestID"

func ContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()

		ctx := context.WithValue(c.Request.Context(), RequestIDKey, requestID)

		c.Request = c.Request.WithContext(ctx)

		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}
